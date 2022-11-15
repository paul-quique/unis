-- Rechercher un utilisateur à partir d'un id de session unique
SELECT user_info.* FROM session JOIN user_info 
ON session.user_id = user_info.id 
WHERE session.id=$1 AND expires_at <= now();

--Ajouter une session dans la base de données, la date d'expiration est générée par le serveur
INSERT INTO session (id, user_id, expires_at) VALUES (:id, :user_id, :expires_at);

--Ajouter un utilisateur dans la base de données, id généré automatiquement
INSERT INTO user_info (first_name, last_name, email, salted_hash, points, created_at)
VALUES (:first_name, :last_name, :email, :salted_hash, :points, :created_at);

--Ajouter une catégorie dans la base de données
INSERT INTO category (name) VALUES (:name);

--Ajouter un produit dans la base de données
INSERT INTO product (name, category_id, price) VALUES (:name, :category_id, :price);

--Récupérer un utilisateur à partir de son identifiant
SELECT * FROM user_info WHERE id=$1;

--Récupérer un utilisateur à partir de son email
SELECT * FROM user_info WHERE email=$1;

--Récupérer une catégorie à partir de son id
SELECT * FROM category WHERE id=$1;

--Récupérer l'ensemble des catégories
SELECT * FROM category;