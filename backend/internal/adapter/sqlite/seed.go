package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

type seedCategory struct {
	nameKK      string
	nameRU      string
	description string
	words       []seedWord
}

type seedWord struct {
	kazakh        string
	russian       string
	transcription string
	exampleKK     string
	exampleRU     string
}

// Seed inserts the starter vocabulary when the database has no categories.
func Seed(ctx context.Context, db *sql.DB) (bool, error) {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories`).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer rollback(tx)

	for _, category := range starter {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO categories (name_kk, name_ru, description)
			VALUES (?, ?, ?)`, category.nameKK, category.nameRU, category.description)
		if err != nil {
			return false, err
		}
		categoryID, err := res.LastInsertId()
		if err != nil {
			return false, err
		}
		for _, word := range category.words {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO words (category_id, kazakh, russian, transcription, example_kk, example_ru)
				VALUES (?, ?, ?, ?, ?, ?)`,
				categoryID, word.kazakh, word.russian, word.transcription, word.exampleKK, word.exampleRU); err != nil {
				return false, fmt.Errorf("seed word %s: %w", word.kazakh, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

var starter = []seedCategory{
	{
		nameKK:      "Сәлемдесу",
		nameRU:      "Приветствия",
		description: "Короткие фразы для начала и конца разговора",
		words: []seedWord{
			{"Сәлем", "Привет", "sälem", "Сәлем!", "Привет!"},
			{"Сәлеметсіз бе", "Здравствуйте", "sälemetsiz be", "Сәлеметсіз бе?", "Здравствуйте?"},
			{"Қайырлы таң", "Доброе утро", "qaiyrly tañ", "Қайырлы таң!", "Доброе утро!"},
			{"Рақмет", "Спасибо", "raqmet", "Рақмет!", "Спасибо!"},
			{"Кешіріңіз", "Извините", "keshiriñiz", "Кешіріңіз.", "Извините."},
			{"Сау болыңыз", "До свидания", "sau bolyñyz", "Сау болыңыз!", "До свидания!"},
		},
	},
	{
		nameKK:      "Отбасы",
		nameRU:      "Семья",
		description: "Родственники",
		words: []seedWord{
			{"Ана", "Мама", "ana", "Ана келді.", "Мама пришла."},
			{"Әке", "Папа", "äke", "Әке үйде.", "Папа дома."},
			{"Аға", "Старший брат", "ağa", "Аға оқиды.", "Старший брат учится."},
			{"Әпке", "Старшая сестра", "äpke", "Әпке ән айтады.", "Старшая сестра поёт."},
			{"Іні", "Младший брат", "ini", "Іні кішкентай.", "Младший брат маленький."},
			{"Әже", "Бабушка", "äje", "Әже шай құяды.", "Бабушка наливает чай."},
		},
	},
	{
		nameKK:      "Тамақ",
		nameRU:      "Еда",
		description: "Еда и напитки",
		words: []seedWord{
			{"Нан", "Хлеб", "nan", "Нан дәмді.", "Хлеб вкусный."},
			{"Су", "Вода", "su", "Су суық.", "Вода холодная."},
			{"Сүт", "Молоко", "süt", "Сүт ақ.", "Молоко белое."},
			{"Алма", "Яблоко", "alma", "Алма қызыл.", "Яблоко красное."},
			{"Ет", "Мясо", "et", "Ет пісті.", "Мясо приготовилось."},
			{"Шай", "Чай", "shai", "Шай ыстық.", "Чай горячий."},
		},
	},
}
