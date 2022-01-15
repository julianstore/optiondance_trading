-- update  instrument name
update `order` set instrument_name = concat("P-pUSD-",instrument_name) where (LENGTH(instrument_name) - LENGTH( REPLACE (instrument_name, '-', ''))) = 3 and id > 0;
update `position` set instrument_name = concat("P-pUSD-",instrument_name) where (LENGTH(instrument_name) - LENGTH( REPLACE (instrument_name, '-', ''))) = 3 and id > 0;
update `trade` set instrument_name = concat("P-pUSD-",instrument_name) where (LENGTH(instrument_name) - LENGTH( REPLACE (instrument_name, '-', ''))) = 3 and id > 0;
update `option_market` set instrument_name = concat("P-pUSD-",instrument_name) where instrument_name != '' and (LENGTH(instrument_name) - LENGTH( REPLACE (instrument_name, '-', ''))) = 3 ;

-- update position option info field
update `position` set delivery_type = (case SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 1), '-', -1) when 'C' then 'CASH' when 'P' then 'PHYSICAL' end)   where id >0;
update `position` set quote_currency = SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 2), '-', -1)   where id > 0;
update `position` set base_currency = SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 3), '-', -1)   where id > 0;
update `position` set strike_price = SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 5), '-', -1)   where id > 0;
update `position` set option_type = (case SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 6), '-', -1) when 'C' then 'CALL' when 'P' then 'PUT' end)   where id > 0;

-- update order delivery type
update `order` set delivery_type = (case SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 1), '-', -1) when 'C' then 'CASH' when 'P' then 'PHYSICAL' end)   where id >0;

-- update option market
update `option_market` set delivery_type = (case SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 1), '-', -1) when 'C' then 'CASH' when 'P' then 'PHYSICAL' end)   where  instrument_name != '';
update `option_market` set quote_currency = SUBSTRING_INDEX(SUBSTRING_INDEX(instrument_name, '-', 2), '-', -1)   where instrument_name != '';