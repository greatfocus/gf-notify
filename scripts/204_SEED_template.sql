DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM template WHERE staticName='client_credentials' THEN
		INSERT INTO template (name, staticName, subject, body, paramsCount)
		VALUES
			('Client Credentials', 'client_credentials', 'Respect Client Credentials', 'Dear Customer,\n\nYou have succefully received your login credentials from your Respect Obituary account.\nYour security question and answer details are as mentioned below.\n\nclient:- $1\nsecret:- $2 \n\nKindly remember the security credentials as they will be needed to reset your password in case you forget your password.\nThank you for choosing our service.\nThank You.\n\n\n\nRegards,\nGreat Focus\n\n\n\n', 1),
		ON CONFLICT
		DO NOTHING;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM template WHERE staticName='email_otp' THEN
		INSERT INTO template (name, staticName, subject, body, paramsCount)
		VALUES
			('Email Otp', 'email_otp', 'Respect OTP for Email Verification', 'Dear Customer,\n\nYou have requested for verification of email address through Respect system.\nThis is to inform you that your request for email verification has been done successfully.\n\nUse following OTP for further proceeding.\n\nOTP For Email Verification:- $1\n\nThank You.\n\n\n\nRegards,\nGreat Focus\n\n\n\n', 1),
		ON CONFLICT
		DO NOTHING;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM template WHERE staticName='first_login' THEN
		INSERT INTO template (name, staticName, subject, body, paramsCount)
		VALUES
			('First Login', 'first_login', 'Respect First Time login : Success', 'Dear Customer,\n\nYou have succefully Completed first time login for your Respect Obituary account.\nKindly remember the security credentials as they will be needed to reset your password in case you forget your password.\nThank you for choosing our service.\nThank You.\n\n\n\nRegards,\nGreat Focus\n\n\n\n', 1),
		ON CONFLICT
		DO NOTHING;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM template WHERE staticName='password_reset' THEN
		INSERT INTO template (name, staticName, subject, body, paramsCount)
		VALUES
			('Password Reset', 'password_reset', 'Respect Password Reset', 'Dear Customer,\n\nPlease note your password has been reset successfully. If you did not make this request, please contact our Customer Experience Center and our team will respond to you soon.\nPlease note our working hours are 0830 to 1630 (EAT) from Monday to Friday and 0900 to 1230 every last Saturday of the month.\nWe regret any delays in reply during non-working hours.\nServing you is our top priority.\n\nThank You.\n\n\n\nRegards,\nGreat Focus\n\n\n\n', 1),
		ON CONFLICT
		DO NOTHING;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM template WHERE staticName='contact_us' THEN
		INSERT INTO template (name, staticName, subject, body, paramsCount)
		VALUES
			('Contact Us', 'contact_us', 'Respect Someone Reached Us', 'Dear Customer,\n\nThank you for contacting Respect Obituary. We have received your e-mail and our team will respond to you soon.\nPlease note our working hours are 0830 to 1630 (EAT) from Monday to Friday and 0900 to 1230 every last Saturday of the month.\nWe regret any delays in reply during non-working hours.\nServing you is our top priority.\n\nThank You.\n\n\n\nRegards,\nGreat Focus\n\n\n\n', 1),
		ON CONFLICT
		DO NOTHING;
	END IF;
END $$;
