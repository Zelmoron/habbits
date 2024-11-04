package controllers

import (
	"Trecker/internal/db"
	"Trecker/internal/db/models"
	"Trecker/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var jwtSecretKey = []byte("very-secret-key")

const (
	contextKeyUser = "user"
)

// Структура HTTP-ответа с информацией о пользователе
type ProfileResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

var (
	errBadCredentials = errors.New("email or password is incorrect")
)

type HabitRequest struct {
	Id int `json:"id"`
}

func Registration(c *fiber.Ctx) error {

	var user UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	db := db.GetDB()

	var userCheck models.UsersModel

	if err := db.Where("email=?", user.Email).First(&userCheck).Error; err == nil {

		return c.Status(409).JSON(fiber.Map{
			"message": "Пользователь уже существует",
		})
	} else if err != gorm.ErrRecordNotFound {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка при проверке пользователя",
		})
	}

	//Если все действия прошли проверки, то хешируем пароль и добавляем в бд
	password, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{
			"mesage": "Ошибка при хешировании пароля",
		})
	}
	var userInsert = models.UsersModel{
		Email:    user.Email,
		Name:     user.Username,
		Password: password,
	}
	if err := db.Create(&userInsert).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка при создании пользователя",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Пользователь добавлен",
	})
}

func Auth(c *fiber.Ctx) error {
	var user UserRequest
	db := db.GetDB()

	// Чтение тела запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	// Логика проверки учетных данных
	var existingUser models.UsersModel
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неправильный логин или пароль",
		})
	}

	// Создание полезной нагрузки для JWT-токена
	payload := jwt.MapClaims{
		"sub": strconv.Itoa(int(existingUser.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Установка срока действия на 10 секунд
	}

	// Создание и подпись токена с использованием алгоритма HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при подписи JWT-токена")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка сервера",
		})
	}

	// Возвращаем токен в ответе
	return c.JSON(LoginResponse{AccessToken: t})
}

// Проверка токена и получение данных профиля
func Profile(c *fiber.Ctx) error {
	token := c.Context().Value(contextKeyUser).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "пользователь не найден",
		})
	}

	// Получаем информацию о пользователе из базы данных
	var user models.UsersModel
	if err := db.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении данных пользователя",
		})
	}

	return c.JSON(ProfileResponse{
		Id:   id,
		Name: user.Name,
	})
}

// Проверка токена и переход на страницу главную

func MainPage(c *fiber.Ctx) error {
	token := c.Context().Value(contextKeyUser).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "пользователь не найден",
		})
	}

	// Получаем информацию о пользователе из базы данных
	var user models.UsersModel
	if err := db.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении данных пользователя",
		})
	}

	return c.JSON(ProfileResponse{
		Id:   id,
		Name: user.Name,
	})
}

func GetHabits(c *fiber.Ctx) error {
	//ПОлучаем привычку
	token := c.Context().Value(contextKeyUser).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "пользователь не найден",
		})
	}

	var habits []models.Habit

	if err := db.GetDB().Where("user_id = ?", id).Find(&habits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON([]models.Habit{})
	}

	return c.JSON(habits)
}

func AddHabits(c *fiber.Ctx) error {
	//Добавляем привычку
	token := c.Context().Value(contextKeyUser).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "пользователь не найден",
		})
	}
	var habit models.Habit

	if err := c.BodyParser(&habit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	habit.UserID = userId // Присваиваем привычке ID пользователя

	if err := db.GetDB().Create(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось добавить привычку"})
	}

	return c.Status(fiber.StatusCreated).JSON(habit)
}

func UpdateHabits(c *fiber.Ctx) error {
	//Обновить день у привычки
	var habit HabitRequest
	// Чтение тела запроса
	if err := c.BodyParser(&habit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	// Получаем токен и проверяем пользователя
	token := c.Context().Value(contextKeyUser).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	_, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "пользователь не найден",
		})
	}

	// Обновляем день у привычки
	if err := db.GetDB().Model(&models.Habit{}).Where("id = ?", habit.Id).Update("day", gorm.Expr("day + 1")).Error; err != nil {
		fmt.Println("texst")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить привычку",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Привычка успешно обновлена",
	})
}
