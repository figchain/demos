from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend

# From client-config.json
private_key_pem = """-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDxwWZP4aIqQ+Fi
WeA6lhRkbMtOeZC2wsabrHs9qjvIexEPi+lh9U6TJWMWUaAcKUIq+4LkHHPJaLJ/
tjc0vwv0kfM8TGimkUAVRk/6gNGlUPwPoqOzh+lLInqMUasGCDjgJOHZ5BA+WJWY
5Qrde9QM1SzfFsMcuyuny+6itMf6cOfT0bMbzxXM0mmoRB7Nsc/MH4dsdLO7PwVg
L/oCRZ7vr2J4IH/6rbX63XjRT0W3nHC23g4fdJCGxUYHqlnDY+NHVFAjAG8uqXED
nORFYfJDwcqzphbqOaQ9oQ7doiRVnTMlfxFR2lw2A+qAHr9YsUkQS85viGPwFc/m
XDOzQ8n1AgMBAAECggEAAI4fpIZRn4Q5Y8buD2Rh1pszWlJIJUtMNnZOcCVQbtTt
hJofJpTwIcFfuWDlm7ryhnO+aLSyBV/irQ1nkgzwQ5cENnq8cMl7mrDnJR2fnaAS
fBy1AIgK3pvNKT5UxLZHHyimhiASc3ozmb7I6BpNMdxZdoewWgQKBSAgQ6pjnS0d
DpyOekXLvPAcV4M67LhtJvFqmrbCH5wPFCmLe3u69Gdc4/a87A5AO5n7DqhLnOW8
HhBJ50Kipb/ShhZv3i7kc77xPbQOzR7QJSgvGv7lGrnFjSm+sxAdQoW/SKjED94A
o6y+8CO8N/r19Wgdc8zasyX6w2b0TyRcIHoAKLFH0QKBgQD9yauGTi7xx5vXKEYm
xBrU+NqikcZhw+0udnSUMZmjgBsZuIlBhqKmMrphuRwFKMFmJ0C/jM6VQ8gszXOi
0BBAiq57ZF603zx3uB3icLQ+ZosMG0X03FzRucDt0undolVCRzVZQ+IZOWDWWAn9
Lok5sui0VpnjKmW3XdW0mvGaOQKBgQDz3OEiCYu16I2ym++d9/uOLwGfvf9k5tmZ
TQeyT+KE/RKCBtBr4XaVj6QYj4VKb6VvQzq0jKEJENxQqbIczlAA6DPi9I4okmur
81W8on+a/pgUm0YAXGzbo1T6E/2kB17kqObehZhlfZuHMwMjRAAWKg1aFgKON139
mHfOtp3dnQKBgAWXf/QBnP8uyrw+4uzPvVeb9BVI6PoWew9fBMqPHTeBBxfV/RA4
izTmQT0N+xQSBdDeZIrT62lWiP6TuEyKERGa/KUzXPLXSFnK8L2ghhgp4j5uC2iN
wm3Mjfevgf+kKATB1OcWm9C6duvCHNY7RELFMmNm1RUwRfV4V9EW6OPhAoGAYQVt
5LbbiOIfDgKQUM8KnEUDZmSXKbPWuvE0sLKrsrFlHapMXb90CIj/hm4DX9wPe7bJ
sm+I2iyFGuqI3IEQv2uiyRb3QBkREXZclBIqpqXIJ9qm/RnIjZHsCxrM/OeZz2uL
fti0CxzwNdgL1YoGZssQSNkc5ywMDwsMD4gEQtkCgYEAoTQTX8HMRpSHr+h99QMz
hpEPZogjE1dEknHS7U5kDfFTIjbQWri19rKyzPD/eRYH4kWmqcbr1ZkdumWBaF8v
UiFc36xW70bTwFowvftlECGo6yHB7g5yyOGcjcW9pjfKfaeEkX55tONqqvLtq9DV
KrrsjwsC1BtY07FGUl9A07Q=
-----END PRIVATE KEY-----"""

# From Cassandra (replacing \n with real newline for key comparison, as cqlsh output escapes it)
public_key_pem_cassandra = """-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8cFmT+GiKkPhYlngOpYU
ZGzLTnmQtsLGm6x7Pao7yHsRD4vpYfVOkyVjFlGgHClCKvuC5BxzyWiyf7Y3NL8L
9JHzPExoppFAFUZP+oDRpVD8D6Kjs4fpSyJ6jFGrBgg44CTh2eQQPliVmOUK3XvU
DNUs3xbDHLsrp8vuorTH+nDn09GzG88VzNJpqEQezbHPzB+HbHSzuz8FYC/6AkWe
769ieCB/+q21+t140U9Ft5xwtt4OH3SQhsVGB6pZw2PjR1RQIwBvLqlxA5zkRWHy
Q8HKs6YW6jmkPaEO3aIkVZ0zJX8RUdpcNgPqgB6/WLFJEEvOb4hj8BXP5lwzs0PJ
9QIDAQAB
-----END PUBLIC KEY-----"""

private_key = serialization.load_pem_private_key(private_key_pem.encode(), password=None, backend=default_backend())
public_key_from_private = private_key.public_key()
public_pem_derived = public_key_from_private.public_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PublicFormat.SubjectPublicKeyInfo
).decode()

# Normalize
def normalize(s):
    return s.strip().replace('\n', '')

if normalize(public_pem_derived) == normalize(public_key_pem_cassandra):
    print("MATCH")
else:
    print("MISMATCH")
    print("Derived (First 50):", normalize(public_pem_derived)[:50])
    print("Cassand (First 50):", normalize(public_key_pem_cassandra)[:50])
