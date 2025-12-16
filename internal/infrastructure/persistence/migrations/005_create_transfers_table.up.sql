CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    seller_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    buyer_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    transfer_price DECIMAL(15,2) NOT NULL,
    transferred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transfers_player_id ON transfers(player_id);
CREATE INDEX idx_transfers_seller_team_id ON transfers(seller_team_id);
CREATE INDEX idx_transfers_buyer_team_id ON transfers(buyer_team_id);

