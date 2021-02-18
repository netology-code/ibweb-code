Для всех конфигураций используется пара user:secret

Генерация файла для Basic аутентификации с алгоритмом bcrypt (по умолчанию - md5)
htpasswd -c -B basic-users.txt user

Генерация файла для Digest аутентификации
htdigest -c digest-users.txt "Digest Authentication" user