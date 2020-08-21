package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	awsprovider "github.com/criticalstack/crit/cmd/bootstrap-server/internal/providers/aws"
	"github.com/criticalstack/crit/pkg/cluster/bootstrap"
	"github.com/criticalstack/crit/pkg/cluster/bootstrap/authorizers/ec2metadata"
)

func (r *bootstrapRouter) handleAmazonIdentityDocumentAndSignature(data []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var sdoc ec2metadata.SignedDocument
		if err := json.Unmarshal(data, &sdoc); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err := ec2metadata.Verify(sdoc.Document, sdoc.Signature); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		var doc ec2metadata.Document
		if err := json.Unmarshal(sdoc.Document, &doc); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cfg := &aws.Config{Region: aws.String(doc.Region)}
		ip, profile, err := awsprovider.GetInstanceInfo(ctx, cfg, doc.InstanceId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if c.RealIP() != ip {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": fmt.Sprintf("expected ip %q, received %q", ip, c.RealIP()),
			})
		}

		for k, v := range r.cfg.Filters {
			switch k {
			case "account-id":
				if doc.AccountId != v {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": fmt.Sprintf("account not authorized: %#v", doc.AccountId),
					})
				}
			case "iam-profile":
				if profile != v {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": fmt.Sprintf("IAM instance profile not authorized: %#v", profile),
					})
				}
			}
		}

		token, err := createNewToken(r.cfg.Kubeconfig)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, &bootstrap.Response{
			BootstrapToken: token,
		})
	}
}
