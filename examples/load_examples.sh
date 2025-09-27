#!/bin/bash

# HaruDB Examples Loader
# This script helps load the example databases for testing

echo "🚀 HaruDB Examples Loader"
echo "========================="
echo ""

# Check if harudb server is running
if ! pgrep -f "harudb" > /dev/null; then
    echo "❌ HaruDB server is not running!"
    echo "Please start the server first: ./harudb --data-dir ./data"
    exit 1
fi

echo "✅ HaruDB server is running"
echo ""

# Function to load example
load_example() {
    local example_name=$1
    local sql_file=$2
    
    echo "📁 Loading $example_name..."
    
    if [ ! -f "$sql_file" ]; then
        echo "❌ File $sql_file not found!"
        return 1
    fi
    
    # Create a temporary file with login commands
    temp_file=$(mktemp)
    echo "LOGIN admin admin123" > "$temp_file"
    cat "$sql_file" >> "$temp_file"
    echo "LOGOUT" >> "$temp_file"
    echo "exit" >> "$temp_file"
    
    # Load the example
    if ../haru-cli < "$temp_file" > /dev/null 2>&1; then
        echo "✅ $example_name loaded successfully!"
    else
        echo "❌ Failed to load $example_name"
        rm "$temp_file"
        return 1
    fi
    
    rm "$temp_file"
    echo ""
}

# Main menu
echo "Select an example to load:"
echo "1) Banking System"
echo "2) Food Ordering App"
echo "3) Both Examples"
echo "4) Exit"
echo ""

read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        load_example "Banking System" "banking_system.sql"
        ;;
    2)
        load_example "Food Ordering App" "food_ordering_app.sql"
        ;;
    3)
        load_example "Banking System" "banking_system.sql"
        load_example "Food Ordering App" "food_ordering_app.sql"
        ;;
    4)
        echo "👋 Goodbye!"
        exit 0
        ;;
    *)
        echo "❌ Invalid choice!"
        exit 1
        ;;
esac

echo "🎉 Examples loaded successfully!"
echo ""
echo "You can now connect to the database and explore the examples:"
echo "  ./haru-cli"
echo "  LOGIN admin admin123"
echo "  SELECT * FROM customers;  # Banking example"
echo "  SELECT * FROM restaurants;  # Food app example"
