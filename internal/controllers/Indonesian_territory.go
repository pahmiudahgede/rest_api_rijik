package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProvinces(c *fiber.Ctx) error {
	provinces, err := services.GetProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve provinces",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Provinces retrieved successfully",
		provinces,
	))
}

// GetRegencies handles the GET request for regencies
func GetRegencies(c *fiber.Ctx) error {
	regencies, err := services.GetRegencies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve regencies",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Regencies retrieved successfully",
		regencies,
	))
}

// GetDistricts handles the GET request for districts
func GetDistricts(c *fiber.Ctx) error {
	districts, err := services.GetDistricts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve districts",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Districts retrieved successfully",
		districts,
	))
}

// GetVillages handles the GET request for villages
func GetVillages(c *fiber.Ctx) error {
	villages, err := services.GetVillages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve villages",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Villages retrieved successfully",
		villages,
	))
}


func GetProvinceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	province, regencies, err := services.GetProvinceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Province not found",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Provinces by id retrieved successfully",
		fiber.Map{
			"id":            province.ID,
			"provinsi_name": province.Name,
			"list_regency":  regencies,
		},
	))
}

func GetRegencyByID(c *fiber.Ctx) error {
	id := c.Params("id")
	regency, districts, err := services.GetRegencyByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Regency not found",
			nil,
		))
	}

	provinces, _ := services.GetProvinces()
	var provinceName string
	for _, province := range provinces {
		if province.ID == regency.ProvinceID {
			provinceName = province.Name
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Regency by id retrieved successfully",
		fiber.Map{
			"id":             regency.ID,
			"province_id":    regency.ProvinceID,
			"province_name":  provinceName,
			"regency_name":   regency.Name,
			"list_districts": districts,
		},
	))
}

func GetDistrictByID(c *fiber.Ctx) error {
	id := c.Params("id")
	district, villages, err := services.GetDistrictByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"District not found",
			nil,
		))
	}

	provinces, _ := services.GetProvinces()
	regencies, _ := services.GetRegencies()
	var provinceName, regencyName string
	for _, province := range provinces {
		if province.ID == district.RegencyID {
			provinceName = province.Name
			break
		}
	}
	for _, regency := range regencies {
		if regency.ID == district.RegencyID {
			regencyName = regency.Name
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"district by id retrieved successfully",
		fiber.Map{
			"id":            district.ID,
			"province_id":   district.RegencyID,
			"regency_id":    district.RegencyID,
			"province_name": provinceName,
			"regency_name":  regencyName,
			"district_name": district.Name,
			"list_villages": villages,
		},
	))
}

func GetVillageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	village, err := services.GetVillageByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Village not found",
			nil,
		))
	}

	provinces, _ := services.GetProvinces()
	regencies, _ := services.GetRegencies()
	districts, _ := services.GetDistricts()

	var provinceName, regencyName, districtName string
	for _, province := range provinces {
		if province.ID == village.ID {
			provinceName = province.Name
		}
	}
	for _, regency := range regencies {
		if regency.ID == village.ID {
			regencyName = regency.Name
		}
	}
	for _, district := range districts {
		if district.ID == village.ID {
			districtName = district.Name
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"villages by id retrieved successfully",
		fiber.Map{
			"id":            village.ID,
			"province_id":   village.ID,
			"regency_id":    village.ID,
			"district_id":   village.ID,
			"province_name": provinceName,
			"regency_name":  regencyName,
			"district_name": districtName,
			"village_name":  village.Name,
		},
	))
}
