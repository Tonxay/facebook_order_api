package middleware

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validatorObj *validator.Validate

func Init() {
	validatorObj = validator.New()

	validatorObj.RegisterValidation("tel", func(fl validator.FieldLevel) bool {
		var pattern = "^(020[25789][0-9]{7}|20[25789][0-9]{7}|[25789][0-9]{7}|030[0-9]{7}|30[0-9]{7})$"
		var telRegex = regexp.MustCompile(pattern)
		return telRegex.MatchString(fl.Field().String())
	})
}

/*
Validate the request body of the request. From tag in struct
*/
func ValidateReqBody(data interface{}) error {

	errs := validatorObj.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			if err != nil {
				// errMsg := fmt.Sprintf("Field '%s' : '%v' | Needs to pass '%s' validation", err.Field(), err.Value(), err.Tag())
				errMsg := fmt.Sprintf("Field '%s' | Needs to pass '%s' validation", err.Field(), err.Tag())
				return errors.New(errMsg)
			}
		}
	}

	return nil
}

/*
Validate the input is a valid UUID
*/
func ValidateUuId(uuid string) error {
	err := validatorObj.Var(uuid, "uuid")
	if err != nil {
		return errors.New("Invalid UUID")
	}

	return nil
}

/*
Validate the input is a valid page and limit values
*/
func ValidatePageAndLimit(page, limit int) error {
	if page < 1 {
		return errors.New("Invalid page value")
	}

	if limit < 1 || limit > 100 {
		return errors.New("Invalid limit value")
	}

	return nil
}

// Parse and validate body
func ParseAndValidateBody(c *fiber.Ctx, out interface{}) error {

	if err := c.BodyParser(out); err != nil {
		return err
	}
	return ValidateReqBody(out)
}

// Parse and validate body
func ParseAndValidateQuery(c *fiber.Ctx, out interface{}) error {
	if err := c.QueryParser(out); err != nil {
		log.Println(err.Error())
		return errors.New("bad_request")
	}
	return ValidateReqBody(out)
}

func ValidatePhone(value int64) bool {
	str := strconv.Itoa(int(value)) // Convert int32 to string

	if len(str) < 2 {
		return false
	}

	prefix := str[:2] // Get first two digits

	return prefix == "20" || prefix == "30"
}

func SetDefaultDateRangeIfEmpty(startDate, endDate string) (string, string) {
	now := time.Now()
	loc := now.Location()

	if startDate == "" {
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		startDate = start.Format("2006-01-02")
	}

	if endDate == "" {
		// last day of current month
		end := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, loc)
		endDate = end.Format("2006-01-02")
	}

	return startDate, endDate
}
