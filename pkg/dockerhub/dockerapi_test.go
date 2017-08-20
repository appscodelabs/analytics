package dockerhub_test

import (
	"testing"

	"github.com/appscode/analytics/pkg/dockerhub"
	"github.com/appscode/log"
)

func TestCollectAnalytics(t *testing.T) {
	mp := make(map[string]string)
	mp["appscode"] = "1zSZGehbu0O62FksNKk5sY5ZeBMQbyz2BFmbUVaf1spQ"
	mp["kubedb"] = "1UexPVL8szwm99T_9ccz__FdOFB_G2g-kpQOW4UkETG0"
	err := dockerhub.CollectAnalytics(mp)
	if err != nil {
		log.Fatalln(err)
	}
}
