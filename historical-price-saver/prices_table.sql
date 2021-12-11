CREATE TABLE prices(
    id STRING(255),
    cryptocurrency STRING(20),
    fiat  STRING(20),
    monthYear DATE,
    priceHighest FLOAT64,
    priceLowest FLOAT64
);
