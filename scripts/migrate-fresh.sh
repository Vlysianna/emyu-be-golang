#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔄 Database Migration Fresh${NC}"
echo ""

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
else
    echo -e "${RED}❌ .env file not found${NC}"
    exit 1
fi

# MySQL connection details
MYSQL_HOST=${DB_HOST:-127.0.0.1}
MYSQL_PORT=${DB_PORT:-3306}
MYSQL_USER=${DB_USER:-root}
MYSQL_PASSWORD=${DB_PASSWORD}
MYSQL_DB=${DB_NAME:-emyu}

# Drop and recreate database
echo -e "${YELLOW}📦 Dropping database...${NC}"
mysql -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASSWORD -e "DROP DATABASE IF EXISTS $MYSQL_DB;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database dropped${NC}"
else
    echo -e "${RED}❌ Failed to drop database${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}🏗️  Creating database...${NC}"
mysql -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database created${NC}"
else
    echo -e "${RED}❌ Failed to create database${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}📋 Running migrations...${NC}"
mysql -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DB < database/schema.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Migrations completed${NC}"
else
    echo -e "${RED}❌ Failed to run migrations${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Migration fresh completed!${NC}"
