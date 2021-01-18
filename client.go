package paprika

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseUrl = "https://www.paprikaapp.com/api/v1/sync/"

type RecipeItem struct {
	Hash string `json:"hash,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type Result struct {
	Result *json.RawMessage
}

type Status struct {
	Recipes    int `json:"recipes,omitempty"`
	Pantry     int `json:"pantry,omitempty"`
	Meals      int `json:"meals,omitempty"`
	Menus      int `json:"menus,omitempty"`
	Groceries  int `json:"groceries,omitempty"`
	Bookmarks  int `json:"bookmarks,omitempty"`
	Menuitems  int `json:"menuitems,omitempty"`
	Categories int `json:"categories,omitempty"`
}

type Bookmark struct {
	URL       string `json:"url,omitempty"`
	Title     string `json:"title,omitempty"`
	UID       string `json:"uid,omitempty"`
	OrderFlag int    `json:"order_flag,omitempty"`
}

type Category struct {
	Name      string `json:"name,omitempty"`
	UID       string `json:"uid,omitempty"`
	ParentUID string `json:"parent_uid,omitempty"`
	OrderFlag int    `json:"order_flag,omitempty"`
}

type GroceryItem struct {
	Aisle      string `json:"aisle,omitempty"`
	UID        string `json:"uid,omitempty"`
	OrderFlag  int    `json:"order_flag,omitempty"`
	Recipe     string `json:"recipe,omitempty"`
	Name       string `json:"name,omitempty"`
	Purchased  bool   `json:"purchased,omitempty"`
	RecipeUID  string `json:"recipe_uid,omitempty"`
	Ingredient string `json:"ingredient,omitempty"`
}

type Meal struct {
	UID       string `json:"uid,omitempty"`
	OrderFlag int    `json:"order_flag,omitempty"`
	RecipeUID string `json:"recipe_uid,omitempty"`
	Date      string `json:"date,omitempty"`
	Type      int    `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
}

type MenuItem struct {
	Name      string `json:"name,omitempty"`
	RecipeUID string `json:"recipe_uid,omitempty"`
	UID       string `json:"uid,omitempty"`
	MenuUID   string `json:"menu_uid,omitempty"`
	OrderFlag int    `json:"order_flag,omitempty"`
}

type PantryItem struct {
	Aisle      string `json:"aisle,omitempty"`
	UID        string `json:"uid,omitempty"`
	Ingredient string `json:"ingredient,omitempty"`
}

type Recipe struct {
	Rating          int      `json:"rating,omitempty"`
	PhotoHash       string   `json:"photo_hash,omitempty"`
	OnFavorites     bool     `json:"on_favorites,omitempty"`
	Photo           string   `json:"photo,omitempty"`
	UID             string   `json:"uid,omitempty"`
	Scale           int      `json:"scale,omitempty"`
	Ingredients     string   `json:"ingredients,omitempty"`
	Source          string   `json:"source,omitempty"`
	Hash            string   `json:"hash,omitempty"`
	SourceURL       string   `json:"source_url,omitempty"`
	Difficulty      string   `json:"difficulty,omitempty"`
	Categories      []string `json:"categories,omitempty"`
	PhotoURL        string   `json:"photo_url,omitempty"`
	CookTime        string   `json:"cook_time,omitempty"`
	Name            string   `json:"name,omitempty"`
	Created         string   `json:"created,omitempty"`
	Notes           string   `json:"notes,omitempty"`
	ImageURL        string   `json:"image_url,omitempty"`
	PrepTime        string   `json:"prep_time,omitempty"`
	Servings        string   `json:"servings,omitempty"`
	NutritionalInfo string   `json:"nutritional_info,omitempty"`
}

type Client struct {
	username   string
	password   string
	httpClient http.Client
}

func NewClient(username, password string) (Client, error) {
	if strings.TrimSpace(username) == "" {
		return Client{}, fmt.Errorf("username must not be empty")
	}

	if strings.TrimSpace(password) == "" {
		return Client{}, fmt.Errorf("password must not be empty")
	}

	return Client{
		httpClient: http.Client{},
		username:   username,
		password:   password,
	}, nil
}

func (c Client) Recipes() ([]RecipeItem, error) {
	rs := []RecipeItem{}
	err := c.get("recipes", &rs)
	if err != nil {
		return nil, err
	}

	return rs, err
}

func (c Client) Recipe(uid string) (Recipe, error) {
	r := Recipe{}
	err := c.get("recipe/"+uid, &r)
	if err != nil {
		return Recipe{}, err
	}

	return r, err
}

func (c Client) Bookmarks() ([]Bookmark, error) {
	return nil, nil
}

func (c Client) prepareGet(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", baseUrl+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %v", path, err)
	}

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}
	req.SetBasicAuth(c.username, c.password)
	return req, nil
}

func (c Client) get(path string, value interface{}) error {
	req, err := c.prepareGet(path)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get %s: %s", path, err)
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s %s", resp.Status, bodyText)
	}

	err = unwrapResult(bodyText, value)
	if err != nil {
		return err
	}

	return nil
}

func unwrapResult(jsonData []byte, value interface{}) error {
	var wrapper Result

	err := json.Unmarshal(jsonData, &wrapper)
	if err != nil {
		return fmt.Errorf("failed to unmarshal result wrapper from %s: %s", string(jsonData), err)
	}
	unwrapped, err := wrapper.Result.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to prepare result for unmarshal: %s", err)
	}
	err = json.Unmarshal(unwrapped, &value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal result from %s: %s", string(unwrapped), err)
	}

	return nil
}
