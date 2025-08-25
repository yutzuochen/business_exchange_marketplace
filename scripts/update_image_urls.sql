-- Update image URLs to use local static files
-- This script updates the first 16 listings to use local images

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/tutoring_center.jpg',
    filename = 'tutoring_center.jpg',
    alt_text = '學園家教中心：教室環境'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 0);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/pet_grooming.jpg',
    filename = 'pet_grooming.jpg',
    alt_text = '寵物美容：美容台與設備'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 1);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/nail_salon.jpg',
    filename = 'nail_salon.jpg',
    alt_text = '髮藝沙龍：造型座位區'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 2);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bakery.jpg',
    filename = 'bakery.jpg',
    alt_text = '烘焙坊：麵包陳列櫃'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 3);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/photo_studio.jpg',
    filename = 'photo_studio.jpg',
    alt_text = '創客空間：工作檯與設備'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 4);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bento_shop.jpg',
    filename = 'bento_shop.jpg',
    alt_text = '便當店：餐盒展示'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 5);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/dessert_shop.jpg',
    filename = 'dessert_shop.jpg',
    alt_text = '豆花店：甜品陳列'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 6);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/toy_store.jpg',
    filename = 'toy_store.jpg',
    alt_text = '玩具店：商品陳列'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 7);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/dry_cleaning.jpg',
    filename = 'dry_cleaning.jpg',
    alt_text = '乾洗店：洗衣設備'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 8);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/scooter_shop.jpg',
    filename = 'scooter_shop.jpg',
    alt_text = '機車行：維修區'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 9);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/stationery_store.jpg',
    filename = 'stationery_store.jpg',
    alt_text = '文具行：商品陳列'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 10);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/yoga_studio.jpg',
    filename = 'yoga_studio.jpg',
    alt_text = '瑜珈教室：練習空間'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 11);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/photo_studio.jpg',
    filename = 'photo_studio.jpg',
    alt_text = '攝影工作室：拍攝空間'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 12);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/hotel_room.jpg',
    filename = 'hotel_room.jpg',
    alt_text = '旅店：客房環境'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 13);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/seafood_market.jpg',
    filename = 'seafood_market.jpg',
    alt_text = '海鮮店：新鮮海產'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 14);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/mountain_cafe.jpg',
    filename = 'mountain_cafe.jpg',
    alt_text = '山谷咖啡：景觀座位'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 15);

-- Add more updates for remaining listings using available images
UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bookstore.jpg',
    filename = 'bookstore.jpg',
    alt_text = '書店：閱讀空間'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 16);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/laundromat.jpg',
    filename = 'laundromat.jpg',
    alt_text = '洗衣店：自助設備'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 17);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/flower_shop.jpg',
    filename = 'flower_shop.jpg',
    alt_text = '花店：花卉陳列'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 18);

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bubble_tea.jpg',
    filename = 'bubble_tea.jpg',
    alt_text = '茶飲店：飲品製作'
WHERE listing_id = (SELECT id FROM listings ORDER BY id LIMIT 1 OFFSET 19);

-- For remaining listings, use generic business images
UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/coffee_shop.jpg',
    filename = 'coffee_shop.jpg'
WHERE listing_id IN (
    SELECT id FROM listings ORDER BY id LIMIT 10 OFFSET 20
);

SELECT 'Image URLs updated successfully!' as status;
