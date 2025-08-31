-- Remove seeded auction status data

DELETE FROM auction_status_ref WHERE status_code IN ('draft', 'active', 'extended', 'closed', 'cancelled');