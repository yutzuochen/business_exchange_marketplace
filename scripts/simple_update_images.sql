-- Simple update to replace external URLs with local ones
UPDATE images SET 
    url = REPLACE(url, 'https://hiyori.cc/wp/wp-content/uploads/2021/03/%E5%AD%B8%E6%A0%A1%E5%92%96%E5%95%A1%E9%A4%A8-Ecole-Cafe7.jpg', 'http://127.0.0.1:8080/static/images/listings/tutoring_center.jpg'),
    filename = 'tutoring_center.jpg'
WHERE url LIKE '%hiyori.cc%';

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/coffee_shop.jpg',
    filename = 'coffee_shop.jpg'
WHERE url LIKE '%worldgymtaiwan%' OR url LIKE '%annieko%' OR url LIKE '%taipei%' OR url LIKE '%hippolife%' OR url LIKE '%fupo%';

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bakery.jpg',
    filename = 'bakery.jpg'
WHERE url LIKE '%chunshuitang%' OR url LIKE '%brotherhotel%';

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/bookstore.jpg',
    filename = 'bookstore.jpg'
WHERE url LIKE '%eslite%';

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/laundromat.jpg',
    filename = 'laundromat.jpg'
WHERE url LIKE '%happyskyblue%';

UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/flower_shop.jpg',
    filename = 'flower_shop.jpg'
WHERE url LIKE '%designwant%';

-- Update all remaining external URLs to use local coffee shop image as fallback
UPDATE images SET 
    url = 'http://127.0.0.1:8080/static/images/listings/coffee_shop.jpg',
    filename = 'coffee_shop.jpg'
WHERE url LIKE 'https://%';

SELECT COUNT(*) as 'Images updated' FROM images WHERE url LIKE 'http://127.0.0.1:8080/static/images/listings/%';
