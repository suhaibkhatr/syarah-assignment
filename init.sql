CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    gift_category VARCHAR(100),
    age_group VARCHAR(50),
    brand VARCHAR(100),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_brand (brand),
    INDEX idx_price (price)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample products
INSERT INTO products (name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at) VALUES
('Teddy Bear', 'Soft and cuddly teddy bear for all ages.', 29.99, 'Toys', 'Romantic', 'All Ages', 'CuddleSoft', true, '2025-07-01 10:00:00', '2025-07-01 10:00:00'),
('Scented Candle Set', 'Lavender and vanilla scented candles for relaxation.', 19.50, 'Home Decor', 'Spa & Wellness', 'Adults', 'AromaGlow', true, '2025-07-01 11:00:00', '2025-07-01 11:00:00'),
('Luxury Chocolate Box', 'Assorted Belgian chocolates in a premium gift box.', 35.00, 'Food', 'Gourmet', 'All Ages', 'ChocoHeaven', false, '2025-07-01 12:00:00', '2025-07-02 08:30:00'),
('Birthday Balloon Bouquet', 'Colorful helium balloons with birthday messages.', 15.75, 'Party Supplies', 'Birthday', 'All Ages', 'BalloonBurst', true, '2025-07-01 13:00:00', '2025-07-01 13:00:00'),
('Personalized Mug', 'Ceramic mug with custom name or message.', 12.00, 'Personalized', 'Office', 'Teens & Adults', 'MugMe', true, '2025-07-01 14:00:00', '2025-07-01 14:00:00'),
('Mini Indoor Plant', 'Cute potted succulent for your desk or shelf.', 9.99, 'Plants', 'Housewarming', 'All Ages', 'GreenJoy', true, '2025-07-01 15:00:00', '2025-07-01 15:00:00'),
('Leather Journal', 'Handcrafted leather-bound notebook for writing or sketching.', 22.49, 'Stationery', 'Creative', 'Teens & Adults', 'WriteCraft', true, '2025-07-01 16:00:00', '2025-07-01 16:00:00'),
('DIY Craft Kit', 'All-in-one craft box for creating handmade gifts.', 18.75, 'Crafts', 'Creative', 'Kids & Teens', 'CraftNest', true, '2025-07-01 17:00:00', '2025-07-01 17:00:00'),
('Couple''s Photo Frame', 'Wooden frame with heart-shaped design for couples.', 14.95, 'Photo Frames', 'Anniversary', 'Adults', 'MemoryLane', false, '2025-07-01 18:00:00', '2025-07-02 09:00:00');