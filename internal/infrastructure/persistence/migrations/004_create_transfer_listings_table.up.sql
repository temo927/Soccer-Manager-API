CREATE TABLE transfer_listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID UNIQUE NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    asking_price DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'sold', 'cancelled')),
    listed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transfer_listings_player_id ON transfer_listings(player_id);
CREATE INDEX idx_transfer_listings_status ON transfer_listings(status);

