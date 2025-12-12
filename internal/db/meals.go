package db

import (
	"context"
	"time"

	"github.com/google/uuid"

	"farm-time/internal/models"
)

// Meal operations

func (db *DB) CreateMeal(ctx context.Context, eventID string, req models.CreateMealRequest) (*models.Meal, error) {
	meal := &models.Meal{
		ID:        uuid.New().String(),
		EventID:   eventID,
		Name:      req.Name,
		MealType:  req.MealType,
		MealDate:  req.MealDate,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		`INSERT INTO meals (id, event_id, name, meal_type, meal_date, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		meal.ID, meal.EventID, meal.Name, meal.MealType, meal.MealDate, meal.Notes, meal.CreatedAt, meal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return meal, nil
}

func (db *DB) GetMeal(ctx context.Context, id string) (*models.Meal, error) {
	var meal models.Meal
	err := db.pool.QueryRow(ctx,
		`SELECT id, event_id, name, meal_type,
		 CASE WHEN meal_date IS NOT NULL THEN meal_date::text ELSE NULL END as meal_date,
		 COALESCE(notes, ''), created_at, updated_at
		 FROM meals WHERE id = $1`, id,
	).Scan(&meal.ID, &meal.EventID, &meal.Name, &meal.MealType, &meal.MealDate, &meal.Notes, &meal.CreatedAt, &meal.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (db *DB) GetMealsByEvent(ctx context.Context, eventID string) ([]models.Meal, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, event_id, name, meal_type,
		 CASE WHEN meal_date IS NOT NULL THEN meal_date::text ELSE NULL END as meal_date,
		 COALESCE(notes, ''), created_at, updated_at
		 FROM meals WHERE event_id = $1 ORDER BY meal_date ASC NULLS LAST, created_at ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var m models.Meal
		if err := rows.Scan(&m.ID, &m.EventID, &m.Name, &m.MealType, &m.MealDate, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}

	return meals, nil
}

func (db *DB) UpdateMeal(ctx context.Context, id string, req models.UpdateMealRequest) (*models.Meal, error) {
	meal, err := db.GetMeal(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		meal.Name = *req.Name
	}
	if req.MealType != nil {
		meal.MealType = *req.MealType
	}
	if req.MealDate != nil {
		meal.MealDate = req.MealDate
	}
	if req.Notes != nil {
		meal.Notes = *req.Notes
	}
	meal.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE meals SET name=$1, meal_type=$2, meal_date=$3, notes=$4, updated_at=$5 WHERE id=$6`,
		meal.Name, meal.MealType, meal.MealDate, meal.Notes, meal.UpdatedAt, id,
	)
	if err != nil {
		return nil, err
	}

	return meal, nil
}

func (db *DB) DeleteMeal(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM meals WHERE id = $1`, id)
	return err
}

// MealItem operations

func (db *DB) CreateMealItem(ctx context.Context, mealID string, req models.CreateMealItemRequest) (*models.MealItem, error) {
	item := &models.MealItem{
		ID:          uuid.New().String(),
		MealID:      mealID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		`INSERT INTO meal_items (id, meal_id, name, description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		item.ID, item.MealID, item.Name, item.Description, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (db *DB) GetMealItem(ctx context.Context, id string) (*models.MealItem, error) {
	var item models.MealItem
	err := db.pool.QueryRow(ctx,
		`SELECT id, meal_id, name, description, created_at, updated_at
		 FROM meal_items WHERE id = $1`, id,
	).Scan(&item.ID, &item.MealID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (db *DB) GetMealItemsByMeal(ctx context.Context, mealID string) ([]models.MealItem, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, meal_id, name, description, created_at, updated_at
		 FROM meal_items WHERE meal_id = $1 ORDER BY name ASC`, mealID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MealItem
	for rows.Next() {
		var i models.MealItem
		if err := rows.Scan(&i.ID, &i.MealID, &i.Name, &i.Description, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (db *DB) UpdateMealItem(ctx context.Context, id string, req models.UpdateMealItemRequest) (*models.MealItem, error) {
	item, err := db.GetMealItem(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	item.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE meal_items SET name=$1, description=$2, updated_at=$3 WHERE id=$4`,
		item.Name, item.Description, item.UpdatedAt, id,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (db *DB) DeleteMealItem(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM meal_items WHERE id = $1`, id)
	return err
}

// MealSignup operations

func (db *DB) CreateMealSignup(ctx context.Context, mealItemID string, userID string, req models.CreateMealSignupRequest) (*models.MealSignup, error) {
	// Get user info for the signup record
	var userName, userEmail string
	err := db.pool.QueryRow(ctx, `SELECT name, email FROM users WHERE id = $1`, userID).Scan(&userName, &userEmail)
	if err != nil {
		return nil, err
	}

	signup := &models.MealSignup{
		ID:         uuid.New().String(),
		MealItemID: mealItemID,
		UserID:     userID,
		UserName:   userName,
		UserEmail:  userEmail,
		Notes:      req.Notes,
		CreatedAt:  time.Now(),
	}

	_, err = db.pool.Exec(ctx,
		`INSERT INTO meal_signups (id, meal_item_id, user_id, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		signup.ID, signup.MealItemID, signup.UserID, signup.Notes, signup.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return signup, nil
}

func (db *DB) GetSignupsByMealItem(ctx context.Context, mealItemID string) ([]models.MealSignup, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT s.id, s.meal_item_id, s.user_id, u.name, u.email, s.notes, s.created_at
		 FROM meal_signups s
		 JOIN users u ON s.user_id = u.id
		 WHERE s.meal_item_id = $1 ORDER BY s.created_at ASC`, mealItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signups []models.MealSignup
	for rows.Next() {
		var s models.MealSignup
		if err := rows.Scan(&s.ID, &s.MealItemID, &s.UserID, &s.UserName, &s.UserEmail, &s.Notes, &s.CreatedAt); err != nil {
			return nil, err
		}
		signups = append(signups, s)
	}

	return signups, nil
}

func (db *DB) DeleteMealSignup(ctx context.Context, mealItemID string, userID string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM meal_signups WHERE meal_item_id = $1 AND user_id = $2`, mealItemID, userID)
	return err
}

// Composite queries

func (db *DB) GetMealItemWithSignups(ctx context.Context, itemID string) (*models.MealItemWithSignups, error) {
	item, err := db.GetMealItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	signups, err := db.GetSignupsByMealItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if signups == nil {
		signups = []models.MealSignup{}
	}

	return &models.MealItemWithSignups{
		MealItem: *item,
		Signups:  signups,
	}, nil
}

func (db *DB) GetMealWithItems(ctx context.Context, mealID string) (*models.MealWithItems, error) {
	meal, err := db.GetMeal(ctx, mealID)
	if err != nil {
		return nil, err
	}

	items, err := db.GetMealItemsByMeal(ctx, mealID)
	if err != nil {
		return nil, err
	}

	// Get signups for each item
	itemsWithSignups := make([]models.MealItemWithSignups, len(items))
	for i, item := range items {
		signups, err := db.GetSignupsByMealItem(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		if signups == nil {
			signups = []models.MealSignup{}
		}
		itemsWithSignups[i] = models.MealItemWithSignups{
			MealItem: item,
			Signups:  signups,
		}
	}

	return &models.MealWithItems{
		Meal:  *meal,
		Items: itemsWithSignups,
	}, nil
}

func (db *DB) GetMealsWithItemsByEvent(ctx context.Context, eventID string) ([]models.MealWithItems, error) {
	meals, err := db.GetMealsByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]models.MealWithItems, len(meals))
	for i, meal := range meals {
		items, err := db.GetMealItemsByMeal(ctx, meal.ID)
		if err != nil {
			return nil, err
		}

		itemsWithSignups := make([]models.MealItemWithSignups, len(items))
		for j, item := range items {
			signups, err := db.GetSignupsByMealItem(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			if signups == nil {
				signups = []models.MealSignup{}
			}
			itemsWithSignups[j] = models.MealItemWithSignups{
				MealItem: item,
				Signups:  signups,
			}
		}

		result[i] = models.MealWithItems{
			Meal:  meal,
			Items: itemsWithSignups,
		}
	}

	return result, nil
}

func (db *DB) GetEventWithMeals(ctx context.Context, id string) (*models.EventWithMeals, error) {
	eventWithAttendees, err := db.GetEventWithAttendees(ctx, id)
	if err != nil {
		return nil, err
	}

	meals, err := db.GetMealsWithItemsByEvent(ctx, id)
	if err != nil {
		return nil, err
	}
	if meals == nil {
		meals = []models.MealWithItems{}
	}

	return &models.EventWithMeals{
		EventWithAttendees: *eventWithAttendees,
		Meals:              meals,
	}, nil
}
