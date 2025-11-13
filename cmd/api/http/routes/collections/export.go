package collections

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

func (r *CollectionsRouter) ExportRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Collections exported successfully.").
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
			Summary:     "Collections Export",
			Description: "Exports the collections available in the system.",
			Tags:        []string{"Collections"},
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
		Path:   "/collections/export",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
			r.middleware.Authorized("collections.export"),
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

			var collections []models.Collection

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
				Find(&collections).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to retrieve collections.",
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
				"Seller Username",
				"Seller ID Number",
				"Material Name",
				"Material GW Code",
				"Material Carbon Factor",
				"Material Weight (kg)",
				"Material Value (ZAR)",
			}}

			for _, collection := range collections {
				collectionRecords := [][]string{}

				collectionDate := collection.CreatedAt.Format("2006-01-02")

				buyerName := collection.Buyer.Name
				buyerAddressParts := []string{}
				buyerBankDetailsParts := []string{}

				if collection.Buyer.Address != nil {
					if collection.Buyer.Address.LineOne != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							collection.Buyer.Address.LineOne,
						)
					}

					if collection.Buyer.Address.LineTwo != nil && *collection.Buyer.Address.LineTwo != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							*collection.Buyer.Address.LineTwo,
						)
					}

					if collection.Buyer.Address.City != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							collection.Buyer.Address.City,
						)
					}

					if collection.Buyer.Address.Province != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							collection.Buyer.Address.Province,
						)
					}

					if collection.Buyer.Address.Country != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							collection.Buyer.Address.Country,
						)
					}

					if collection.Buyer.Address.ZipCode != "" {
						buyerAddressParts = append(
							buyerAddressParts,
							collection.Buyer.Address.ZipCode,
						)
					}
				}

				if collection.Buyer.BankDetails != nil {
					if collection.Buyer.BankDetails.AccountHolder != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Account Holder: "+collection.Buyer.BankDetails.AccountHolder,
						)
					}

					if collection.Buyer.BankDetails.AccountNumber != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Account Number: "+collection.Buyer.BankDetails.AccountNumber,
						)
					}

					if collection.Buyer.BankDetails.BankName != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Bank Name: "+collection.Buyer.BankDetails.BankName,
						)
					}

					if collection.Buyer.BankDetails.BranchCode != "" {
						buyerBankDetailsParts = append(
							buyerBankDetailsParts,
							"Branch Code: "+collection.Buyer.BankDetails.BranchCode,
						)
					}
				}

				sellerName := collection.Seller.Name
				sellerAddressParts := []string{}
				sellerBankDetailsParts := []string{}
				sellerUsername := collection.Seller.Username
				sellerIdNumber := "Not Provided"

				if collection.Seller.IdNumber != nil {
					sellerIdNumber = *collection.Seller.IdNumber
				}

				if collection.Seller.Address != nil {
					if collection.Seller.Address.LineOne != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							collection.Seller.Address.LineOne,
						)
					}

					if collection.Seller.Address.LineTwo != nil && *collection.Seller.Address.LineTwo != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							*collection.Seller.Address.LineTwo,
						)
					}

					if collection.Seller.Address.City != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							collection.Seller.Address.City,
						)
					}

					if collection.Seller.Address.Province != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							collection.Seller.Address.Province,
						)
					}

					if collection.Seller.Address.Country != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							collection.Seller.Address.Country,
						)
					}

					if collection.Seller.Address.ZipCode != "" {
						sellerAddressParts = append(
							sellerAddressParts,
							collection.Seller.Address.ZipCode,
						)
					}
				}

				if collection.Seller.BankDetails != nil {
					if collection.Seller.BankDetails.AccountHolder != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Account Holder: "+collection.Seller.BankDetails.AccountHolder,
						)
					}

					if collection.Seller.BankDetails.AccountNumber != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Account Number: "+collection.Seller.BankDetails.AccountNumber,
						)
					}

					if collection.Seller.BankDetails.BankName != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Bank Name: "+collection.Seller.BankDetails.BankName,
						)
					}

					if collection.Seller.BankDetails.BranchCode != "" {
						sellerBankDetailsParts = append(
							sellerBankDetailsParts,
							"Branch Code: "+collection.Seller.BankDetails.BranchCode,
						)
					}
				}

				collectionFirstRecord := []string{
					collectionDate,
					buyerName,
					strings.Join(buyerAddressParts, ", "),
					strings.Join(buyerBankDetailsParts, ", "),
					sellerName,
					strings.Join(sellerAddressParts, ", "),
					strings.Join(sellerBankDetailsParts, ", "),
					sellerUsername,
					sellerIdNumber,
				}

				if len(collection.Materials) > 0 {
					totalCollectionValue := 0.0

					for materialIdx, material := range collection.Materials {
						if materialIdx == 0 {
							collectionFirstRecord = append(
								collectionFirstRecord,
								material.Name,
								material.GWCode,
								material.CarbonFactor,
								fmt.Sprintf("%.2f", material.Weight),
								fmt.Sprintf("%.2f", material.Value),
							)

							totalCollectionValue += material.Value * material.Weight

							collectionRecords = append(
								collectionRecords,
								collectionFirstRecord,
							)
						} else {
							materialRecord := []string{
								"", "", "", "", "", "", "", "", "",
								material.Name,
								material.GWCode,
								material.CarbonFactor,
								fmt.Sprintf("%.2f", material.Weight),
								fmt.Sprintf("%.2f", material.Value),
							}

							collectionRecords = append(
								collectionRecords,
								materialRecord,
							)

							totalCollectionValue += material.Value * material.Weight
						}
					}

					collectionRecords = append(
						collectionRecords,
						[]string{
							"", "", "", "", "", "", "", "", "", "", "", "", "Total:",
							fmt.Sprintf("%.2f", totalCollectionValue),
						},
					)
				}

				records = append(
					records,
					collectionRecords...,
				)
			}

			c.Set("Content-Type", "text/csv")
			c.Set(
				"Content-Disposition",
				fmt.Sprintf(
					"attachment; filename=collections_export_%s.csv",
					time.Now().Format(time.TimeOnly),
				),
			)

			csvWriter := csv.NewWriter(c.Response().BodyWriter())

			csvWriter.WriteAll(records)

			return nil
		},
	}
}
