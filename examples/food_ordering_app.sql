-- Food Ordering App Example for HaruDB (Swiggy/Zomato style)
-- This example demonstrates a complete food delivery system with restaurants, menus, orders, and delivery tracking

-- =============================================
-- FOOD ORDERING APP SCHEMA
-- =============================================

-- 1. Users table - stores customer information
CREATE TABLE users (
    user_id, 
    username, 
    email, 
    phone, 
    password_hash, 
    full_name, 
    address, 
    city, 
    state, 
    zip_code, 
    created_at, 
    status
);

-- 2. Restaurants table - stores restaurant information
CREATE TABLE restaurants (
    restaurant_id, 
    name, 
    description, 
    cuisine_type, 
    address, 
    city, 
    state, 
    zip_code, 
    phone, 
    email, 
    rating, 
    delivery_fee, 
    minimum_order, 
    delivery_time, 
    created_at, 
    status
);

-- 3. Menu_items table - stores restaurant menu items
CREATE TABLE menu_items (
    item_id, 
    restaurant_id, 
    name, 
    description, 
    price, 
    category, 
    is_vegetarian, 
    is_vegan, 
    is_gluten_free, 
    calories, 
    preparation_time, 
    created_at, 
    status
);

-- 4. Orders table - stores order information
CREATE TABLE orders (
    order_id, 
    user_id, 
    restaurant_id, 
    order_date, 
    total_amount, 
    delivery_fee, 
    tax_amount, 
    final_amount, 
    delivery_address, 
    special_instructions, 
    payment_method, 
    order_status, 
    created_at
);

-- 5. Order_items table - stores individual items in each order
CREATE TABLE order_items (
    order_item_id, 
    order_id, 
    item_id, 
    quantity, 
    unit_price, 
    total_price, 
    special_instructions
);

-- 6. Delivery_persons table - stores delivery person information
CREATE TABLE delivery_persons (
    delivery_person_id, 
    name, 
    phone, 
    email, 
    vehicle_type, 
    license_number, 
    rating, 
    status, 
    created_at
);

-- 7. Deliveries table - stores delivery tracking information
CREATE TABLE deliveries (
    delivery_id, 
    order_id, 
    delivery_person_id, 
    pickup_time, 
    delivery_time, 
    delivery_status, 
    delivery_address, 
    notes, 
    created_at
);

-- 8. Reviews table - stores customer reviews
CREATE TABLE reviews (
    review_id, 
    user_id, 
    restaurant_id, 
    order_id, 
    rating, 
    comment, 
    review_date, 
    created_at
);

-- 9. Categories table - stores food categories
CREATE TABLE categories (
    category_id, 
    name, 
    description, 
    icon_url
);

-- 10. Promotions table - stores promotional offers
CREATE TABLE promotions (
    promotion_id, 
    restaurant_id, 
    title, 
    description, 
    discount_type, 
    discount_value, 
    minimum_order, 
    valid_from, 
    valid_until, 
    status
);

-- =============================================
-- CREATE INDEXES FOR PERFORMANCE
-- =============================================

CREATE INDEX ON users (email);
CREATE INDEX ON users (phone);
CREATE INDEX ON restaurants (city);
CREATE INDEX ON restaurants (cuisine_type);
CREATE INDEX ON menu_items (restaurant_id);
CREATE INDEX ON menu_items (category);
CREATE INDEX ON orders (user_id);
CREATE INDEX ON orders (restaurant_id);
CREATE INDEX ON orders (order_date);
CREATE INDEX ON order_items (order_id);
CREATE INDEX ON deliveries (order_id);
CREATE INDEX ON reviews (restaurant_id);

-- =============================================
-- SAMPLE DATA INSERTION
-- =============================================

-- Insert sample users
INSERT INTO users VALUES (1, 'john_doe', 'john.doe@email.com', '555-1001', 'hashed_password_1', 'John Doe', '123 Main St', 'New York', 'NY', '10001', '2023-01-15', 'active');
INSERT INTO users VALUES (2, 'sarah_smith', 'sarah.smith@email.com', '555-1002', 'hashed_password_2', 'Sarah Smith', '456 Oak Ave', 'Los Angeles', 'CA', '90210', '2023-02-20', 'active');
INSERT INTO users VALUES (3, 'mike_wilson', 'mike.wilson@email.com', '555-1003', 'hashed_password_3', 'Mike Wilson', '789 Pine St', 'Chicago', 'IL', '60601', '2023-03-10', 'active');
INSERT INTO users VALUES (4, 'emily_davis', 'emily.davis@email.com', '555-1004', 'hashed_password_4', 'Emily Davis', '321 Elm St', 'Houston', 'TX', '77001', '2023-04-05', 'active');
INSERT INTO users VALUES (5, 'david_brown', 'david.brown@email.com', '555-1005', 'hashed_password_5', 'David Brown', '654 Maple Ave', 'Phoenix', 'AZ', '85001', '2023-05-12', 'active');

-- Insert sample categories
INSERT INTO categories VALUES (1, 'Pizza', 'Italian pizza and related dishes', 'pizza_icon.png');
INSERT INTO categories VALUES (2, 'Burgers', 'Burgers, sandwiches, and wraps', 'burger_icon.png');
INSERT INTO categories VALUES (3, 'Asian', 'Chinese, Japanese, Thai, and other Asian cuisines', 'asian_icon.png');
INSERT INTO categories VALUES (4, 'Indian', 'Traditional Indian cuisine', 'indian_icon.png');
INSERT INTO categories VALUES (5, 'Mexican', 'Mexican and Tex-Mex cuisine', 'mexican_icon.png');
INSERT INTO categories VALUES (6, 'Desserts', 'Cakes, ice cream, and sweet treats', 'dessert_icon.png');

-- Insert sample restaurants
INSERT INTO restaurants VALUES (1, 'Mario''s Pizza Palace', 'Authentic Italian pizza made with fresh ingredients', 'Italian', '100 Pizza St', 'New York', 'NY', '10001', '555-2001', 'info@mariospizza.com', 4.5, 2.99, 15.00, 30, '2023-01-01', 'active');
INSERT INTO restaurants VALUES (2, 'Burger King Deluxe', 'Premium burgers and American classics', 'American', '200 Burger Ave', 'Los Angeles', 'CA', '90210', '555-2002', 'contact@burgerkingdeluxe.com', 4.2, 3.50, 12.00, 25, '2023-01-15', 'active');
INSERT INTO restaurants VALUES (3, 'Golden Dragon', 'Authentic Chinese cuisine with modern twists', 'Chinese', '300 Dragon St', 'Chicago', 'IL', '60601', '555-2003', 'hello@goldendragon.com', 4.7, 2.50, 20.00, 35, '2023-02-01', 'active');
INSERT INTO restaurants VALUES (4, 'Spice Palace', 'Traditional Indian dishes with authentic flavors', 'Indian', '400 Spice Rd', 'Houston', 'TX', '77001', '555-2004', 'info@spicepalace.com', 4.6, 3.00, 18.00, 40, '2023-02-15', 'active');
INSERT INTO restaurants VALUES (5, 'Taco Fiesta', 'Fresh Mexican food and vibrant flavors', 'Mexican', '500 Fiesta Blvd', 'Phoenix', 'AZ', '85001', '555-2005', 'hola@tacofiesta.com', 4.3, 2.75, 10.00, 20, '2023-03-01', 'active');

-- Insert sample menu items
INSERT INTO menu_items VALUES (1, 1, 'Margherita Pizza', 'Classic tomato, mozzarella, and basil', 14.99, 'Pizza', 'true', 'false', 'false', 280, 15, '2023-01-01', 'active');
INSERT INTO menu_items VALUES (2, 1, 'Pepperoni Pizza', 'Pepperoni with mozzarella cheese', 16.99, 'Pizza', 'false', 'false', 'false', 320, 15, '2023-01-01', 'active');
INSERT INTO menu_items VALUES (3, 1, 'Veggie Supreme', 'Bell peppers, onions, mushrooms, olives', 15.99, 'Pizza', 'true', 'true', 'true', 250, 15, '2023-01-01', 'active');
INSERT INTO menu_items VALUES (4, 2, 'Classic Cheeseburger', 'Beef patty with cheese, lettuce, tomato', 12.99, 'Burgers', 'false', 'false', 'false', 450, 10, '2023-01-15', 'active');
INSERT INTO menu_items VALUES (5, 2, 'Chicken Deluxe', 'Grilled chicken with avocado and bacon', 14.99, 'Burgers', 'false', 'false', 'false', 380, 12, '2023-01-15', 'active');
INSERT INTO menu_items VALUES (6, 3, 'Kung Pao Chicken', 'Spicy chicken with peanuts and vegetables', 13.99, 'Asian', 'false', 'false', 'true', 320, 20, '2023-02-01', 'active');
INSERT INTO menu_items VALUES (7, 3, 'Vegetable Lo Mein', 'Stir-fried noodles with mixed vegetables', 11.99, 'Asian', 'true', 'true', 'false', 280, 15, '2023-02-01', 'active');
INSERT INTO menu_items VALUES (8, 4, 'Chicken Curry', 'Spicy chicken curry with basmati rice', 16.99, 'Indian', 'false', 'false', 'true', 420, 25, '2023-02-15', 'active');
INSERT INTO menu_items VALUES (9, 4, 'Vegetable Biryani', 'Fragrant rice with mixed vegetables', 14.99, 'Indian', 'true', 'true', 'true', 350, 20, '2023-02-15', 'active');
INSERT INTO menu_items VALUES (10, 5, 'Beef Tacos', 'Three soft tacos with seasoned beef', 9.99, 'Mexican', 'false', 'false', 'true', 280, 8, '2023-03-01', 'active');

-- Insert sample delivery persons
INSERT INTO delivery_persons VALUES (1, 'Alex Johnson', '555-3001', 'alex.johnson@delivery.com', 'Bicycle', 'BIKE001', 4.8, 'active', '2023-01-01');
INSERT INTO delivery_persons VALUES (2, 'Maria Garcia', '555-3002', 'maria.garcia@delivery.com', 'Motorcycle', 'MOTO001', 4.6, 'active', '2023-01-15');
INSERT INTO delivery_persons VALUES (3, 'Tom Wilson', '555-3003', 'tom.wilson@delivery.com', 'Car', 'CAR001', 4.9, 'active', '2023-02-01');
INSERT INTO delivery_persons VALUES (4, 'Lisa Chen', '555-3004', 'lisa.chen@delivery.com', 'Bicycle', 'BIKE002', 4.7, 'active', '2023-02-15');
INSERT INTO delivery_persons VALUES (5, 'James Brown', '555-3005', 'james.brown@delivery.com', 'Motorcycle', 'MOTO002', 4.5, 'active', '2023-03-01');

-- Insert sample orders
INSERT INTO orders VALUES (1, 1, 1, '2023-06-15 18:30:00', 32.97, 2.99, 2.64, 38.60, '123 Main St, New York, NY 10001', 'Please ring doorbell', 'credit_card', 'delivered', '2023-06-15 18:30:00');
INSERT INTO orders VALUES (2, 2, 2, '2023-06-16 12:15:00', 27.98, 3.50, 2.24, 33.72, '456 Oak Ave, Los Angeles, CA 90210', 'Leave at door', 'debit_card', 'delivered', '2023-06-16 12:15:00');
INSERT INTO orders VALUES (3, 3, 3, '2023-06-17 19:45:00', 25.98, 2.50, 2.08, 30.56, '789 Pine St, Chicago, IL 60601', 'Call when arrived', 'credit_card', 'delivered', '2023-06-17 19:45:00');
INSERT INTO orders VALUES (4, 4, 4, '2023-06-18 20:00:00', 31.98, 3.00, 2.56, 37.54, '321 Elm St, Houston, TX 77001', 'Extra spicy please', 'cash', 'delivered', '2023-06-18 20:00:00');
INSERT INTO orders VALUES (5, 5, 5, '2023-06-19 13:20:00', 19.98, 2.75, 1.60, 24.33, '654 Maple Ave, Phoenix, AZ 85001', 'No onions', 'credit_card', 'delivered', '2023-06-19 13:20:00');

-- Insert sample order items
INSERT INTO order_items VALUES (1, 1, 1, 1, 14.99, 14.99, 'Extra cheese');
INSERT INTO order_items VALUES (2, 1, 2, 1, 16.99, 16.99, 'Well done');
INSERT INTO order_items VALUES (3, 2, 4, 1, 12.99, 12.99, 'No pickles');
INSERT INTO order_items VALUES (4, 2, 5, 1, 14.99, 14.99, 'Extra bacon');
INSERT INTO order_items VALUES (5, 3, 6, 1, 13.99, 13.99, 'Medium spicy');
INSERT INTO order_items VALUES (6, 3, 7, 1, 11.99, 11.99, 'Extra vegetables');
INSERT INTO order_items VALUES (7, 4, 8, 1, 16.99, 16.99, 'Extra spicy');
INSERT INTO order_items VALUES (8, 4, 9, 1, 14.99, 14.99, 'Mild spice');
INSERT INTO order_items VALUES (9, 5, 10, 2, 9.99, 19.98, 'No onions');

-- Insert sample deliveries
INSERT INTO deliveries VALUES (1, 1, 1, '2023-06-15 18:45:00', '2023-06-15 19:15:00', 'delivered', '123 Main St, New York, NY 10001', 'Delivered successfully', '2023-06-15 18:45:00');
INSERT INTO deliveries VALUES (2, 2, 2, '2023-06-16 12:30:00', '2023-06-16 12:55:00', 'delivered', '456 Oak Ave, Los Angeles, CA 90210', 'Left at door as requested', '2023-06-16 12:30:00');
INSERT INTO deliveries VALUES (3, 3, 3, '2023-06-17 20:00:00', '2023-06-17 20:35:00', 'delivered', '789 Pine St, Chicago, IL 60601', 'Called customer before delivery', '2023-06-17 20:00:00');
INSERT INTO deliveries VALUES (4, 4, 4, '2023-06-18 20:15:00', '2023-06-18 20:55:00', 'delivered', '321 Elm St, Houston, TX 77001', 'Extra spicy order delivered', '2023-06-18 20:15:00');
INSERT INTO deliveries VALUES (5, 5, 5, '2023-06-19 13:35:00', '2023-06-19 13:55:00', 'delivered', '654 Maple Ave, Phoenix, AZ 85001', 'No onions as requested', '2023-06-19 13:35:00');

-- Insert sample reviews
INSERT INTO reviews VALUES (1, 1, 1, 1, 5, 'Amazing pizza! Fresh ingredients and perfect crust.', '2023-06-16', '2023-06-16 10:30:00');
INSERT INTO reviews VALUES (2, 2, 2, 2, 4, 'Great burger, but delivery was a bit slow.', '2023-06-17', '2023-06-17 09:15:00');
INSERT INTO reviews VALUES (3, 3, 3, 3, 5, 'Excellent Chinese food! Will order again.', '2023-06-18', '2023-06-18 11:45:00');
INSERT INTO reviews VALUES (4, 4, 4, 4, 5, 'Authentic Indian flavors, highly recommended!', '2023-06-19', '2023-06-19 14:20:00');
INSERT INTO reviews VALUES (5, 5, 5, 5, 4, 'Good tacos, fresh ingredients.', '2023-06-20', '2023-06-20 16:10:00');

-- Insert sample promotions
INSERT INTO promotions VALUES (1, 1, '20% Off First Order', 'Get 20% off your first pizza order', 'percentage', 20, 20.00, '2023-06-01', '2023-12-31', 'active');
INSERT INTO promotions VALUES (2, 2, 'Free Delivery', 'Free delivery on orders over $25', 'fixed', 3.50, 25.00, '2023-06-01', '2023-12-31', 'active');
INSERT INTO promotions VALUES (3, 3, 'Combo Deal', 'Buy 2 entrees get 1 free', 'fixed', 13.99, 30.00, '2023-06-01', '2023-12-31', 'active');
INSERT INTO promotions VALUES (4, 4, 'Weekend Special', '15% off all orders on weekends', 'percentage', 15, 20.00, '2023-06-01', '2023-12-31', 'active');
INSERT INTO promotions VALUES (5, 5, 'Taco Tuesday', 'Buy 2 tacos get 1 free on Tuesdays', 'fixed', 9.99, 20.00, '2023-06-01', '2023-12-31', 'active');

-- =============================================
-- SAMPLE QUERIES FOR FOOD ORDERING APP
-- =============================================

-- Query 1: Get all restaurants with their ratings and delivery info
SELECT name, cuisine_type, rating, delivery_fee, minimum_order, delivery_time 
FROM restaurants 
WHERE status = 'active' 
ORDER BY rating DESC;

-- Query 2: Get menu items for a specific restaurant
SELECT name, description, price, category, is_vegetarian, is_vegan 
FROM menu_items 
WHERE restaurant_id = 1 AND status = 'active' 
ORDER BY category, price;

-- Query 3: Get order history for a specific user
SELECT o.order_id, r.name as restaurant_name, o.order_date, o.final_amount, o.order_status 
FROM orders o, restaurants r 
WHERE o.user_id = 1 AND o.restaurant_id = r.restaurant_id 
ORDER BY o.order_date DESC;

-- Query 4: Get order details with items
SELECT o.order_id, r.name as restaurant_name, mi.name as item_name, oi.quantity, oi.total_price 
FROM orders o, restaurants r, order_items oi, menu_items mi 
WHERE o.order_id = 1 AND o.restaurant_id = r.restaurant_id AND oi.order_id = o.order_id AND mi.item_id = oi.item_id;

-- Query 5: Get delivery tracking information
SELECT o.order_id, r.name as restaurant_name, dp.name as delivery_person, d.delivery_status, d.pickup_time, d.delivery_time 
FROM orders o, restaurants r, deliveries d, delivery_persons dp 
WHERE o.order_id = 1 AND o.restaurant_id = r.restaurant_id AND d.order_id = o.order_id AND dp.delivery_person_id = d.delivery_person_id;

-- Query 6: Get restaurant reviews and ratings
SELECT r.name as restaurant_name, u.username, rev.rating, rev.comment, rev.review_date 
FROM restaurants r, reviews rev, users u 
WHERE r.restaurant_id = 1 AND rev.restaurant_id = r.restaurant_id AND u.user_id = rev.user_id 
ORDER BY rev.review_date DESC;

-- Query 7: Get popular menu items (most ordered)
SELECT mi.name, mi.price, COUNT(oi.order_item_id) as order_count, SUM(oi.quantity) as total_quantity 
FROM menu_items mi, order_items oi 
WHERE mi.item_id = oi.item_id 
GROUP BY mi.item_id, mi.name, mi.price 
ORDER BY total_quantity DESC;

-- Query 8: Get delivery person performance
SELECT dp.name, dp.vehicle_type, dp.rating, COUNT(d.delivery_id) as total_deliveries 
FROM delivery_persons dp, deliveries d 
WHERE dp.delivery_person_id = d.delivery_person_id AND d.delivery_status = 'delivered' 
GROUP BY dp.delivery_person_id, dp.name, dp.vehicle_type, dp.rating 
ORDER BY dp.rating DESC;

-- Query 9: Get active promotions
SELECT p.title, r.name as restaurant_name, p.description, p.discount_type, p.discount_value, p.minimum_order 
FROM promotions p, restaurants r 
WHERE p.restaurant_id = r.restaurant_id AND p.status = 'active' 
AND p.valid_from <= NOW() AND p.valid_until >= NOW();

-- Query 10: Get user order statistics
SELECT u.username, COUNT(o.order_id) as total_orders, SUM(o.final_amount) as total_spent, AVG(o.final_amount) as avg_order_value 
FROM users u, orders o 
WHERE u.user_id = o.user_id 
GROUP BY u.user_id, u.username 
ORDER BY total_spent DESC;
