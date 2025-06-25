#!/bin/bash

# Create a simple KTM icon using ImageMagick or sips
# This creates a basic text icon - you can customize colors and style

echo "Creating KTM icon..."

# Check if ImageMagick is installed
if command -v convert &> /dev/null; then
    echo "Using ImageMagick to create icon..."
    
    # Create icon sizes
    sizes=(16 32 64 128 256 512 1024)
    
    mkdir -p icon.iconset
    
    for size in "${sizes[@]}"; do
        echo "Creating ${size}x${size} icon..."
        convert -size ${size}x${size} \
                -background "#2C3E50" \
                -fill white \
                -gravity center \
                -font "Helvetica-Bold" \
                -pointsize $((size/3)) \
                label:"KTM" \
                "icon.iconset/icon_${size}x${size}.png"
    done
    
    # Create @2x versions
    echo "Creating icon_16x16@2x.png..."
    cp icon.iconset/icon_32x32.png icon.iconset/icon_16x16@2x.png
    echo "Creating icon_32x32@2x.png..."
    cp icon.iconset/icon_64x64.png icon.iconset/icon_32x32@2x.png
    echo "Creating icon_128x128@2x.png..."
    cp icon.iconset/icon_256x256.png icon.iconset/icon_128x128@2x.png
    echo "Creating icon_256x256@2x.png..."
    cp icon.iconset/icon_512x512.png icon.iconset/icon_256x256@2x.png
    echo "Creating icon_512x512@2x.png..."
    cp icon.iconset/icon_1024x1024.png icon.iconset/icon_512x512@2x.png
    
    # Remove the 1024 file (we only need it as 512@2x)
    rm icon.iconset/icon_1024x1024.png
    
    # Generate icns file
    echo "Generating icon.icns..."
    iconutil -c icns icon.iconset
    
    # Clean up
    rm -rf icon.iconset
    
    echo "✅ Icon created: icon.icns"
    
else
    echo "ImageMagick not found. Creating a placeholder icon using macOS tools..."
    
    # Create a simple icon using sips and other built-in tools
    mkdir -p icon.iconset
    
    # Create a basic 1024x1024 image with text
    # This is a fallback method using system tools
    cat > create_icon.py << 'EOF'
#!/usr/bin/env python3
import os
from PIL import Image, ImageDraw, ImageFont
import sys

# Create icon sizes
sizes = [16, 32, 64, 128, 256, 512, 1024]

os.makedirs('icon.iconset', exist_ok=True)

for size in sizes:
    # Create image with dark blue background
    img = Image.new('RGB', (size, size), color='#2C3E50')
    draw = ImageDraw.Draw(img)
    
    # Try to use a system font
    try:
        font_size = int(size * 0.4)
        font = ImageFont.truetype('/System/Library/Fonts/Helvetica.ttc', font_size)
    except:
        font = ImageFont.load_default()
    
    # Draw text
    text = "KTM"
    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    position = ((size - text_width) // 2, (size - text_height) // 2)
    draw.text(position, text, font=font, fill='white')
    
    # Save
    img.save(f'icon.iconset/icon_{size}x{size}.png')

# Create @2x versions
import shutil
shutil.copy('icon.iconset/icon_32x32.png', 'icon.iconset/icon_16x16@2x.png')
shutil.copy('icon.iconset/icon_64x64.png', 'icon.iconset/icon_32x32@2x.png')
shutil.copy('icon.iconset/icon_256x256.png', 'icon.iconset/icon_128x128@2x.png')
shutil.copy('icon.iconset/icon_512x512.png', 'icon.iconset/icon_256x256@2x.png')
shutil.copy('icon.iconset/icon_1024x1024.png', 'icon.iconset/icon_512x512@2x.png')
os.remove('icon.iconset/icon_1024x1024.png')

print("Icon images created")
EOF

    # Check if Python PIL is available
    if python3 -c "import PIL" 2>/dev/null; then
        python3 create_icon.py
        rm create_icon.py
        
        # Generate icns
        iconutil -c icns icon.iconset
        rm -rf icon.iconset
        echo "✅ Icon created: icon.icns"
    else
        echo "❌ Neither ImageMagick nor Python PIL found."
        echo "Please install one of them:"
        echo "  brew install imagemagick"
        echo "  OR"
        echo "  pip3 install Pillow"
    fi
fi