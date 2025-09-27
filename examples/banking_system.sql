-- Banking System Example for HaruDB
-- This example demonstrates a complete banking system with customers, accounts, transactions, and loans

-- =============================================
-- BANKING SYSTEM SCHEMA
-- =============================================

-- 1. Customers table - stores customer information
CREATE TABLE customers (
    customer_id, 
    first_name, 
    last_name, 
    email, 
    phone, 
    date_of_birth, 
    address, 
    city, 
    state, 
    zip_code, 
    created_at, 
    status
);

-- 2. Accounts table - stores bank accounts
CREATE TABLE accounts (
    account_id, 
    customer_id, 
    account_type, 
    account_number, 
    balance, 
    interest_rate, 
    created_at, 
    status
);

-- 3. Transactions table - stores all financial transactions
CREATE TABLE transactions (
    transaction_id, 
    account_id, 
    transaction_type, 
    amount, 
    description, 
    reference_number, 
    created_at, 
    status
);

-- 4. Loans table - stores loan information
CREATE TABLE loans (
    loan_id, 
    customer_id, 
    loan_type, 
    principal_amount, 
    interest_rate, 
    term_months, 
    monthly_payment, 
    remaining_balance, 
    created_at, 
    status
);

-- 5. Loan_payments table - stores loan payment history
CREATE TABLE loan_payments (
    payment_id, 
    loan_id, 
    payment_amount, 
    principal_payment, 
    interest_payment, 
    payment_date, 
    status
);

-- 6. Branches table - stores bank branch information
CREATE TABLE branches (
    branch_id, 
    branch_name, 
    address, 
    city, 
    state, 
    zip_code, 
    phone, 
    manager_name
);

-- 7. Employees table - stores bank employee information
CREATE TABLE employees (
    employee_id, 
    first_name, 
    last_name, 
    email, 
    phone, 
    position, 
    branch_id, 
    salary, 
    hire_date, 
    status
);

-- =============================================
-- CREATE INDEXES FOR PERFORMANCE
-- =============================================

CREATE INDEX ON customers (email);
CREATE INDEX ON customers (phone);
CREATE INDEX ON accounts (customer_id);
CREATE INDEX ON accounts (account_number);
CREATE INDEX ON transactions (account_id);
CREATE INDEX ON transactions (created_at);
CREATE INDEX ON loans (customer_id);
CREATE INDEX ON loan_payments (loan_id);
CREATE INDEX ON employees (branch_id);

-- =============================================
-- SAMPLE DATA INSERTION
-- =============================================

-- Insert sample customers
INSERT INTO customers VALUES (1, 'John', 'Smith', 'john.smith@email.com', '555-0101', '1985-03-15', '123 Main St', 'New York', 'NY', '10001', '2023-01-15', 'active');
INSERT INTO customers VALUES (2, 'Sarah', 'Johnson', 'sarah.johnson@email.com', '555-0102', '1990-07-22', '456 Oak Ave', 'Los Angeles', 'CA', '90210', '2023-02-20', 'active');
INSERT INTO customers VALUES (3, 'Michael', 'Brown', 'michael.brown@email.com', '555-0103', '1988-11-08', '789 Pine St', 'Chicago', 'IL', '60601', '2023-03-10', 'active');
INSERT INTO customers VALUES (4, 'Emily', 'Davis', 'emily.davis@email.com', '555-0104', '1992-05-14', '321 Elm St', 'Houston', 'TX', '77001', '2023-04-05', 'active');
INSERT INTO customers VALUES (5, 'David', 'Wilson', 'david.wilson@email.com', '555-0105', '1987-09-30', '654 Maple Ave', 'Phoenix', 'AZ', '85001', '2023-05-12', 'active');

-- Insert sample branches
INSERT INTO branches VALUES (1, 'Downtown Branch', '100 Financial Plaza', 'New York', 'NY', '10001', '555-1001', 'Robert Manager');
INSERT INTO branches VALUES (2, 'Westside Branch', '200 Commerce St', 'Los Angeles', 'CA', '90210', '555-1002', 'Lisa Director');
INSERT INTO branches VALUES (3, 'Central Branch', '300 Business Ave', 'Chicago', 'IL', '60601', '555-1003', 'James Supervisor');

-- Insert sample employees
INSERT INTO employees VALUES (1, 'Robert', 'Manager', 'robert.manager@bank.com', '555-2001', 'Branch Manager', 1, 75000, '2020-01-15', 'active');
INSERT INTO employees VALUES (2, 'Lisa', 'Director', 'lisa.director@bank.com', '555-2002', 'Branch Director', 2, 85000, '2019-06-01', 'active');
INSERT INTO employees VALUES (3, 'James', 'Supervisor', 'james.supervisor@bank.com', '555-2003', 'Branch Supervisor', 3, 65000, '2021-03-10', 'active');
INSERT INTO employees VALUES (4, 'Maria', 'Teller', 'maria.teller@bank.com', '555-2004', 'Bank Teller', 1, 45000, '2022-08-20', 'active');
INSERT INTO employees VALUES (5, 'Tom', 'Advisor', 'tom.advisor@bank.com', '555-2005', 'Financial Advisor', 2, 60000, '2021-11-15', 'active');

-- Insert sample accounts
INSERT INTO accounts VALUES (1, 1, 'checking', 'CHK001234567', 5000.00, 0.01, '2023-01-15', 'active');
INSERT INTO accounts VALUES (2, 1, 'savings', 'SAV001234567', 15000.00, 2.50, '2023-01-15', 'active');
INSERT INTO accounts VALUES (3, 2, 'checking', 'CHK002345678', 3200.00, 0.01, '2023-02-20', 'active');
INSERT INTO accounts VALUES (4, 2, 'savings', 'SAV002345678', 8500.00, 2.50, '2023-02-20', 'active');
INSERT INTO accounts VALUES (5, 3, 'checking', 'CHK003456789', 1200.00, 0.01, '2023-03-10', 'active');
INSERT INTO accounts VALUES (6, 4, 'checking', 'CHK004567890', 7500.00, 0.01, '2023-04-05', 'active');
INSERT INTO accounts VALUES (7, 4, 'savings', 'SAV004567890', 25000.00, 2.50, '2023-04-05', 'active');
INSERT INTO accounts VALUES (8, 5, 'checking', 'CHK005678901', 2800.00, 0.01, '2023-05-12', 'active');

-- Insert sample transactions
INSERT INTO transactions VALUES (1, 1, 'deposit', 5000.00, 'Initial deposit', 'DEP001', '2023-01-15 09:00:00', 'completed');
INSERT INTO transactions VALUES (2, 2, 'deposit', 15000.00, 'Initial deposit', 'DEP002', '2023-01-15 09:05:00', 'completed');
INSERT INTO transactions VALUES (3, 1, 'withdrawal', 500.00, 'ATM withdrawal', 'WTH001', '2023-01-20 14:30:00', 'completed');
INSERT INTO transactions VALUES (4, 1, 'transfer', 1000.00, 'Transfer to savings', 'TRF001', '2023-01-25 10:15:00', 'completed');
INSERT INTO transactions VALUES (5, 2, 'deposit', 1000.00, 'Transfer from checking', 'TRF001', '2023-01-25 10:15:00', 'completed');
INSERT INTO transactions VALUES (6, 3, 'deposit', 3200.00, 'Initial deposit', 'DEP003', '2023-02-20 11:00:00', 'completed');
INSERT INTO transactions VALUES (7, 4, 'deposit', 8500.00, 'Initial deposit', 'DEP004', '2023-02-20 11:05:00', 'completed');
INSERT INTO transactions VALUES (8, 1, 'withdrawal', 200.00, 'ATM withdrawal', 'WTH002', '2023-02-28 16:45:00', 'completed');
INSERT INTO transactions VALUES (9, 6, 'deposit', 7500.00, 'Initial deposit', 'DEP005', '2023-04-05 13:20:00', 'completed');
INSERT INTO transactions VALUES (10, 7, 'deposit', 25000.00, 'Initial deposit', 'DEP006', '2023-04-05 13:25:00', 'completed');

-- Insert sample loans
INSERT INTO loans VALUES (1, 1, 'personal', 10000.00, 8.50, 36, 315.00, 10000.00, '2023-06-01', 'active');
INSERT INTO loans VALUES (2, 2, 'auto', 25000.00, 6.25, 60, 485.00, 25000.00, '2023-07-15', 'active');
INSERT INTO loans VALUES (3, 4, 'home', 300000.00, 4.75, 360, 1564.00, 300000.00, '2023-08-10', 'active');
INSERT INTO loans VALUES (4, 5, 'personal', 5000.00, 9.00, 24, 228.00, 5000.00, '2023-09-05', 'active');

-- Insert sample loan payments
INSERT INTO loan_payments VALUES (1, 1, 315.00, 250.00, 65.00, '2023-07-01', 'completed');
INSERT INTO loan_payments VALUES (2, 1, 315.00, 252.00, 63.00, '2023-08-01', 'completed');
INSERT INTO loan_payments VALUES (3, 2, 485.00, 350.00, 135.00, '2023-08-15', 'completed');
INSERT INTO loan_payments VALUES (4, 3, 1564.00, 1200.00, 364.00, '2023-09-10', 'completed');
INSERT INTO loan_payments VALUES (5, 4, 228.00, 190.00, 38.00, '2023-10-05', 'completed');

-- =============================================
-- SAMPLE QUERIES FOR BANKING SYSTEM
-- =============================================

-- Query 1: Get all customers with their account balances
SELECT c.first_name, c.last_name, c.email, a.account_type, a.balance 
FROM customers c, accounts a 
WHERE c.customer_id = a.customer_id;

-- Query 2: Get transaction history for a specific account
SELECT transaction_type, amount, description, created_at 
FROM transactions 
WHERE account_id = 1 
ORDER BY created_at DESC;

-- Query 3: Get customers with high balance accounts (>$10000)
SELECT c.first_name, c.last_name, a.account_type, a.balance 
FROM customers c, accounts a 
WHERE c.customer_id = a.customer_id AND a.balance > 10000;

-- Query 4: Get loan information with customer details
SELECT c.first_name, c.last_name, l.loan_type, l.principal_amount, l.interest_rate, l.monthly_payment 
FROM customers c, loans l 
WHERE c.customer_id = l.customer_id;

-- Query 5: Get total deposits and withdrawals for each account
SELECT account_id, 
       SUM(CASE WHEN transaction_type = 'deposit' THEN amount ELSE 0 END) as total_deposits,
       SUM(CASE WHEN transaction_type = 'withdrawal' THEN amount ELSE 0 END) as total_withdrawals
FROM transactions 
GROUP BY account_id;

-- Query 6: Get employees by branch
SELECT b.branch_name, e.first_name, e.last_name, e.position, e.salary 
FROM branches b, employees e 
WHERE b.branch_id = e.branch_id 
ORDER BY b.branch_name, e.salary DESC;

-- Query 7: Get loan payment history
SELECT c.first_name, c.last_name, l.loan_type, lp.payment_amount, lp.payment_date 
FROM customers c, loans l, loan_payments lp 
WHERE c.customer_id = l.customer_id AND l.loan_id = lp.loan_id 
ORDER BY lp.payment_date DESC;
