package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupIncomeRoutesPart(route fiber.Router) {
	// 1. ຂໍ້ກຳນົດ Middleware (JWT)
	// ບັນທັດນີ້ຈະເຮັດໃຫ້ທຸກ route ຂ້າງລຸ່ມຕ້ອງມີ Token
	route.Use(middleware.JWTProtected)

	// ==========================================
	// Income Handlers (ລະບຸ Path ເຕັມໆ ທີ່ນີ້ເລີຍ)
	// ==========================================

	// ສ້າງລາຍຮັບ (Batch) -> POST /incomes
	route.Post("/incomes", service.CreateIncomesInBatch)

	// ດຶງລາຍຮັບລາຍເດືອນ -> GET /incomes
	route.Get("/incomes", service.GetIncomesHandler)

	// ==========================================
	// Category Handlers
	// ==========================================

	// ດຶງຫມວດຫມູ່ -> GET /incomes/categories
	route.Get("/incomes/categories", service.GetIncomeCategoriesHandler)

	// ສ້າງຫມວດຫມູ່ -> POST /incomes/categories
	route.Post("/incomes/categories", service.CreateIncomeCategoryHandler)

	// ==========================================
	// Report Handlers
	// ==========================================

	// ເບິ່ງຍອດສະຫຼຸບ (Dashboard) -> GET /reports/dashboard
	route.Get("/reports/dashboard", service.GetFinancialStatusHandler)

	// ເບິ່ງ Statement -> GET /reports/statement
	route.Get("/reports/statement", service.GetStatementHandler)
}
