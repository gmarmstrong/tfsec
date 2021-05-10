package checks

import (
	"fmt"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"

	"github.com/tfsec/tfsec/internal/app/tfsec/parser"
)

const AWSUnencryptedAtRestElasticacheReplicationGroup scanner.RuleCode = "AWS035"
const AWSUnencryptedAtRestElasticacheReplicationGroupDescription scanner.RuleSummary = "Unencrypted Elasticache Replication Group."
const AWSUnencryptedAtRestElasticacheReplicationGroupImpact = "Data in the replication group could be readable if compromised"
const AWSUnencryptedAtRestElasticacheReplicationGroupResolution = "Enable encryption for replication group"
const AWSUnencryptedAtRestElasticacheReplicationGroupExplanation = `
You should ensure your Elasticache data is encrypted at rest to help prevent sensitive information from being read by unauthorised users.
`
const AWSUnencryptedAtRestElasticacheReplicationGroupBadExample = `
resource "aws_elasticache_replication_group" "bad_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        at_rest_encryption_enabled = false
}
`
const AWSUnencryptedAtRestElasticacheReplicationGroupGoodExample = `
resource "aws_elasticache_replication_group" "good_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        at_rest_encryption_enabled = true
}
`

func init() {
	scanner.RegisterCheck(scanner.Check{
		Code: AWSUnencryptedAtRestElasticacheReplicationGroup,
		Documentation: scanner.CheckDocumentation{
			Summary:     AWSUnencryptedAtRestElasticacheReplicationGroupDescription,
			Impact:      AWSUnencryptedAtRestElasticacheReplicationGroupImpact,
			Resolution:  AWSUnencryptedAtRestElasticacheReplicationGroupResolution,
			Explanation: AWSUnencryptedAtRestElasticacheReplicationGroupExplanation,
			BadExample:  AWSUnencryptedAtRestElasticacheReplicationGroupBadExample,
			GoodExample: AWSUnencryptedAtRestElasticacheReplicationGroupGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_replication_group#at_rest_encryption_enabled",
				"https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/at-rest-encryption.html",
			},
		},
		Provider:       scanner.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticache_replication_group"},
		CheckFunc: func(check *scanner.Check, block *parser.Block, context *scanner.Context) []scanner.Result {

			encryptionAttr := block.GetAttribute("at_rest_encryption_enabled")
			if encryptionAttr == nil {
				return []scanner.Result{
					check.NewResult(
						fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (missing at_rest_encryption_enabled attribute).", block.FullName()),
						block.Range(),
						scanner.SeverityError,
					),
				}
			} else if !isBooleanOrStringTrue(encryptionAttr) {
				return []scanner.Result{
					check.NewResultWithValueAnnotation(
						fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (at_rest_encryption_enabled set to false).", block.FullName()),
						encryptionAttr.Range(),
						encryptionAttr,
						scanner.SeverityError,
					),
				}
			}

			return nil
		},
	})
}
