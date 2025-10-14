package transactions

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *TransactionsRouter) ExportRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Transactions exported successfully.").
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchemaRef(schemas.SuccessSchema),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Transactions Export",
			Description: "Exports the transactions available in the system.",
			Tags:        []string{"Transactions"},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewQueryParameter("startDate").
						WithRequired(true).
						WithSchema(openapi3.NewDateTimeSchema().WithFormat("yyyy-MM-dd")),
				},
				{
					Value: openapi3.NewQueryParameter("endDate").
						WithRequired(true).
						WithSchema(openapi3.NewDateTimeSchema().WithFormat("yyyy-MM-dd")),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path:   "/transactions/export",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
			r.middleware.Authorized("transactions.export"),
		},
		Handler: func(c *fiber.Ctx) error {
			startDate := c.Query(
				"startDate",
				time.Now().Format("2006-01-02"),
			)
			endDate := c.Query(
				"endDate",
				time.Now().Format("2006-01-02"),
			)

			var transactions []models.Transaction

			if err := r.storage.
				Database().
				Preload("Seller").
				Preload("Buyer").
				Preload("Materials").
				Where(
					"created_at >= ? AND created_at <= ?",
					startDate,
					endDate,
				).
				Find(&transactions).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to retrieve transactions.",
				})
			}

			records := [][]string{{
				"Date",
				"Buyer Name",
				"Buyer Address",
				"Buyer Bank Details",
				"Seller Name",
				"Seller Address",
				"Seller Bank Details",
				"Material Name",
				"Material GW Code",
				"Material Carbon Factor",
				"Material Weight (kg)",
				"Material Value (ZAR)",
			}}

			for _, transaction := range transactions {
				transactionRecords := [][]string{}

				transactionDate := transaction.CreatedAt.Format("2006-01-02")

				buyerName := transaction.Buyer.Name
				buyerAddressParts := []string{}
				buyerBankDetailsParts := []string{}

				if transaction.Buyer.Address != nil {
					if transaction.Buyer.Address.LineOne != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							transaction.Buyer.Address.LineOne,
						)
					}

					if transaction.Buyer.Address.LineTwo != nil && *transaction.Buyer.Address.LineTwo != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							*transaction.Buyer.Address.LineTwo,
						)
					}

					if transaction.Buyer.Address.City != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							transaction.Buyer.Address.City,
						)
					}

					if transaction.Buyer.Address.Province != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							transaction.Buyer.Address.Province,
						)
					}

					if transaction.Buyer.Address.Country != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							transaction.Buyer.Address.Country,
						)
					}

					if transaction.Buyer.Address.ZipCode != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							transaction.Buyer.Address.ZipCode,
						)
					}
				}

				if transaction.Buyer.BankDetails != nil {
					if transaction.Buyer.BankDetails.AccountHolder != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Account Holder: "+transaction.Buyer.BankDetails.AccountHolder,
						)
					}

					if transaction.Buyer.BankDetails.AccountNumber != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Account Number: "+transaction.Buyer.BankDetails.AccountNumber,
						)
					}

					if transaction.Buyer.BankDetails.BankName != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Bank Name: "+transaction.Buyer.BankDetails.BankName,
						)
					}

					if transaction.Buyer.BankDetails.BranchCode != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Branch Code: "+transaction.Buyer.BankDetails.BranchCode,
						)
					}
				}

				sellerName := transaction.Seller.Name
				sellerAddressParts := []string{}
				sellerBankDetailsParts := []string{}

				if transaction.Seller.Address != nil {
					if transaction.Seller.Address.LineOne != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							transaction.Seller.Address.LineOne,
						)
					}

					if transaction.Seller.Address.LineTwo != nil && *transaction.Seller.Address.LineTwo != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							*transaction.Seller.Address.LineTwo,
						)
					}

					if transaction.Seller.Address.City != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							transaction.Seller.Address.City,
						)
					}

					if transaction.Seller.Address.Province != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							transaction.Seller.Address.Province,
						)
					}

					if transaction.Seller.Address.Country != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							transaction.Seller.Address.Country,
						)
					}

					if transaction.Seller.Address.ZipCode != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							transaction.Seller.Address.ZipCode,
						)
					}
				}

				if transaction.Seller.BankDetails != nil {
					if transaction.Seller.BankDetails.AccountHolder != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Account Holder: "+transaction.Seller.BankDetails.AccountHolder,
						)
					}

					if transaction.Seller.BankDetails.AccountNumber != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Account Number: "+transaction.Seller.BankDetails.AccountNumber,
						)
					}

					if transaction.Seller.BankDetails.BankName != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Bank Name: "+transaction.Seller.BankDetails.BankName,
						)
					}

					if transaction.Seller.BankDetails.BranchCode != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Branch Code: "+transaction.Seller.BankDetails.BranchCode,
						)
					}
				}

				transactionFirstRecord := []string{
					transactionDate,
					buyerName,
					strings.Join(buyerAddressParts, ", "),
					strings.Join(buyerBankDetailsParts, ", "),
					sellerName,
					strings.Join(sellerAddressParts, ", "),
					strings.Join(sellerBankDetailsParts, ", "),
				}

				if len(transaction.Materials) > 0 {
					totalTransactionValue := 0.0

					for materialIdx, material := range transaction.Materials {
						if materialIdx == 0 {
							transactionFirstRecord = append(
								transactionFirstRecord,
								material.Name,
								material.GWCode,
								material.CarbonFactor,
								fmt.Sprintf("%.2f", material.Weight),
								fmt.Sprintf("%.2f", material.Value),
							)

							totalTransactionValue += material.Value * material.Weight

							transactionRecords = append(
								transactionRecords,
								transactionFirstRecord,
							)
						} else {
							materialRecord := []string{
								"", "", "", "", "", "", "",
								material.Name,
								material.GWCode,
								material.CarbonFactor,
								fmt.Sprintf("%.2f", material.Weight),
								fmt.Sprintf("%.2f", material.Value),
							}

							transactionRecords = append(
								transactionRecords,
								materialRecord,
							)

							totalTransactionValue += material.Value * material.Weight
						}
					}

					transactionRecords = append(
						transactionRecords,
						[]string{
							"", "", "", "", "", "", "", "", "", "", "Total:",
							fmt.Sprintf("%.2f", totalTransactionValue),
						},
					)
				}

				records = append(
					records,
					transactionRecords...,
				)
			}

			c.Set("Content-Type", "text/csv")
			c.Set(
				"Content-Disposition",
				fmt.Sprintf(
					"attachment; filename=transactions_export_%s.csv",
					time.Now().Format(time.TimeOnly),
				),
			)

			csvWriter := csv.NewWriter(c.Response().BodyWriter())

			csvWriter.WriteAll(records)

			return nil
		},
	}
}
