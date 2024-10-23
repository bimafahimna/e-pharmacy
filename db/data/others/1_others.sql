-- CUSTOMER DETAILS
insert into "customer_details" ("user_id","profile_photo_url") values
('1','https://res.cloudinary.com/dgdauj2jq/image/upload/v1701356937/test_photo.jpg'),
('2','https://res.cloudinary.com/dgdauj2jq/image/upload/v1701356937/test_photo.jpg'),
('3','https://res.cloudinary.com/dgdauj2jq/image/upload/v1701356937/test_photo.jpg');

-- MANUFACTURER
insert into "manufacturers" ("name") values
('PT Kimia Farma Tbk'),
('PT Kalbe Farma Tbk');

-- PRODUCT CLASSIFICATIONS
insert into "product_classifications" ("name") values
('Obat Bebas'),
('Obat Keras'),
('Obat Bebas Terbatas'),
('Non Obat');

-- PRODUCT FORMS
insert into "product_forms" ("name") values
('Tablet'),
('Kapsul'),
('Pil'),
('Serbuk'),
('Supositoria'),
('Ovula'),
('Salep'),
('Krim'),
('Gel'),
('Sirup'),
('Suspensi'),
('Eliksir'),
('Infus'),
('Tetes'),
('Inhalasi'),
('Aerosol'),
('Turbuhaler');

-- CUSTOMER ADDRESSES
insert into "customer_addresses" ("user_id","name","receiver_name","receiver_phone_number","latitude","longitude","province","city_id","city","district","sub_district","address_details","is_active") values
('1','My lovely home','Joni','+6283219923','6.1903','106.7913','JAWA BARAT','23','BANDUNG','Coblong','Dago','Jalan ke rumah tercozy','true'),
('2','My big bad home','BUDI','+62810010021','6.2990','106.6139','JAWA BARAT','78','BOGOR','Bogor Selatan','Bojongkerta','Jalan ke rumah-rumahan','true');

-- CARTS
insert into "carts" ("user_id","pharmacy_id", "product_id","quantity") values
-- customer1@gmail.com
('1', '1', '1', '2'), -- Panadol from pharmacy a, 2 pcs
('1', '2', '3', '1'), -- Tylenol from pharmacy b, 1 pcs
('1', '1', '8', '1'), -- Dulcolax from pharmacy a, 1 pcs
-- customer2@gmail.com
('2', '1', '2', '3'), -- Paramex from pharmacy a, 3 pcs
('2', '2', '5', '2'), -- Bodrex from pharmacy b, 2 pcs
('2', '2', '7', '1'), -- Vicks VapoRub from pharmacy b, 1 pcs
-- customer3@gmail.com
('3', '1', '4', '2'), -- Decolgen from pharmacy a, 2 pcs
('3', '2', '9', '1'), -- Mylanta from pharmacy b, 1 pcs
('3', '2', '10', '2'); -- Neozep from pharmacy b, 2 pcs
