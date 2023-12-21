--- Data used to seed the orders table
INSERT INTO public.orders (id, category, product_type, product_name, stock, product_id, shipping_address, customer_email)
VALUES
    (1, 'Electronics', 'Laptop', 'Example Laptop 1', true, 101, '123 Main St, Cityville', 'customer1@example.com'),
    (2, 'Clothing', 'Shoes', 'Running Shoes 1', true, 102, '456 Oak St, Townsville', 'customer2@example.com'),
    (3, 'Home Goods', 'Furniture', 'Coffee Table 1', true, 103, '789 Pine St, Villageton', 'customer3@example.com'),
    (4, 'Electronics', 'Smartphone', 'Example Phone 1', true, 104, '101 Elm St, Hamletown', 'customer4@example.com'),
    (5, 'Clothing', 'T-Shirt', 'Graphic Tee 1', true, 105, '202 Birch St, Woodsville', 'customer5@example.com'),
    (6, 'Home Goods', 'Appliance', 'Microwave 1', true, 106, '303 Maple St, Orchard City', 'customer6@example.com'),
    (7, 'Electronics', 'Headphones', 'Over-ear Headphones 1', true, 107, '404 Cedar St, Forestville', 'customer7@example.com'),
    (8, 'Clothing', 'Jeans', 'Denim Jeans 1', true, 108, '505 Redwood St, Groveton', 'customer8@example.com'),
    (9, 'Home Goods', 'Lighting', 'Desk Lamp 1', true, 109, '606 Walnut St, Riverside', 'customer9@example.com'),
    (10, 'Electronics', 'Tablet', 'Example Tablet 1', true, 110, '707 Pineapple St, Beachville', 'customer10@example.com');