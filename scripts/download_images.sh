#!/bin/bash

# Create images directory if it doesn't exist
mkdir -p static/images/listings

echo "📸 Downloading sample business listing images..."

# Download sample images for different business types
declare -A images=(
    ["coffee_shop.jpg"]="https://images.unsplash.com/photo-1554118811-1e0d58224f24?w=800"
    ["fitness_gym.jpg"]="https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800"
    ["dessert_shop.jpg"]="https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=800"
    ["kindergarten.jpg"]="https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?w=800"
    ["nail_salon.jpg"]="https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=800"
    ["claw_machine.jpg"]="https://images.unsplash.com/photo-1511512578047-dfb367046420?w=800"
    ["bubble_tea.jpg"]="https://images.unsplash.com/photo-1525385133512-2f3bdd039054?w=800"
    ["bento_shop.jpg"]="https://images.unsplash.com/photo-1563379091339-03246963d51a?w=800"
    ["bookstore.jpg"]="https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=800"
    ["laundromat.jpg"]="https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800"
    ["flower_shop.jpg"]="https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=800"
    ["yoga_studio.jpg"]="https://images.unsplash.com/photo-1506629905645-b178a0ee8c87?w=800"
    ["photo_studio.jpg"]="https://images.unsplash.com/photo-1606983340126-99ab4feaa64a?w=800"
    ["hotel_room.jpg"]="https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?w=800"
    ["seafood_market.jpg"]="https://images.unsplash.com/photo-1544943910-4c1dc44aab44?w=800"
    ["mountain_cafe.jpg"]="https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=800"
    ["stationery_store.jpg"]="https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=800"
    ["scooter_shop.jpg"]="https://images.unsplash.com/photo-1558618047-3c8c76ca7d13?w=800"
    ["dry_cleaning.jpg"]="https://images.unsplash.com/photo-1582735689369-4fe89db7114c?w=800"
    ["toy_store.jpg"]="https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=800"
    ["bakery.jpg"]="https://images.unsplash.com/photo-1509440159596-0249088772ff?w=800"
    ["pet_grooming.jpg"]="https://images.unsplash.com/photo-1601758228041-f3b2795255f1?w=800"
    ["car_wash.jpg"]="https://images.unsplash.com/photo-1563720223185-11003d516935?w=800"
    ["eyeglass_store.jpg"]="https://images.unsplash.com/photo-1574258495973-f010dfbb5371?w=800"
    ["hotpot_restaurant.jpg"]="https://images.unsplash.com/photo-1555396273-367ea4eb4db5?w=800"
    ["tutoring_center.jpg"]="https://images.unsplash.com/photo-1497486751825-1233686d5d80?w=800"
)

# Download each image
for filename in "${!images[@]}"; do
    url="${images[$filename]}"
    echo "📥 Downloading $filename..."
    
    if curl -L -o "static/images/listings/$filename" "$url" --silent --show-error; then
        echo "✅ Downloaded $filename"
    else
        echo "❌ Failed to download $filename"
    fi
    
    # Small delay to be respectful
    sleep 1
done

echo "🎉 Image download complete!"
echo "📁 Images saved to: static/images/listings/"
echo "🌐 Access via: http://localhost:8080/static/images/listings/[filename]"
