package service

import (
	"backend/common/producer"
	"backend/notification/internal/repository"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	_ "github.com/joho/godotenv/autoload"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type Service struct {
	repository    *repository.Repository
	producer      *producer.Producer
	signer        *sign.URLSigner
	cloudfrontURL string
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	pk, err := sign.LoadPEMPrivKeyFile("cert/aws/aws-cloudfront-private-key.pem")
	if err != nil {
		log.Panicf("fail to make cloud front private key: %v", err)
	}
	signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_ID"), pk)

	s := Service{
		repository:    r,
		producer:      p,
		signer:        signer,
		cloudfrontURL: os.Getenv("AWS_CLOUDFRONT_URL"),
	}

	go s.electConversationNotificationLeader()

	return &s
}

func (s *Service) generateSignedURL(contentType, filename string) (string, error) {
	signedURL, err := s.signer.Sign(
		fmt.Sprintf("%s/%s/%s",
			s.cloudfrontURL, contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL",
			"err", err)
		return "", err
	}
	slog.Info("success to sign url", "signedURL", signedURL)
	return signedURL, nil
}

func (s *Service) electConversationNotificationLeader() {
	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	podName := os.Getenv("HOSTNAME")
	namespace := "backend"
	leaseName := "conversation-notification"

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseName,
			Namespace: namespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podName,
		},
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,

		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				slog.Info("pod won the election, starting cron...", "podName", podName)
				s.executeConversationNotifications(ctx)
			},
			OnStoppedLeading: func() {
				slog.Info("pod give up leadership, Shutting down...", "podName", podName)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				if identity != podName {
					slog.Info("pod lose election, stand by...", "currentLeader", identity)
				}
			},
		},
	})
}
