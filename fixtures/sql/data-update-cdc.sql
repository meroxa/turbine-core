UPDATE public.orders SET
    category = 'Electronics Updated',
    shipping_address = 'Updated Address 1'
WHERE id = 1;

UPDATE public.orders SET
    product_type = 'Shoes Updated',
    customer_email = 'updated_customer2@example.com'
WHERE id = 2;

DELETE FROM public.orders WHERE id = 3;

UPDATE public.orders SET
    category = 'Electronics Updated',
    product_type = 'Smartphone Updated',
    shipping_address = 'Updated Address 4'
WHERE id = 4;

DELETE FROM public.orders WHERE id = 5;

UPDATE public.orders SET
    shipping_address = 'Updated Address 6',
    customer_email = 'updated_customer6@example.com'
WHERE id = 6;

UPDATE public.orders SET
    category = 'Electronics Updated',
    product_type = 'Headphones Updated',
    stock = false
WHERE id = 7;

DELETE FROM public.orders WHERE id = 8;

UPDATE public.orders SET
    stock = false,
    customer_email = 'updated_customer9@example.com'
WHERE id = 9;

DELETE FROM public.orders WHERE id = 10;
