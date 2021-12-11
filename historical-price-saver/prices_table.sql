CREATE TABLE prices (
  id STRING(255),
  cryptocurrency STRING(20),
  fiat STRING(20),
  priceHighest FLOAT64,
  priceLowest FLOAT64,
  monthYear TIMESTAMP,
) PRIMARY KEY(id);
