# SIGame

## Содержание
- [Русский](#русский)
    - [Краткое описание](#краткое-описание)
    - [Запуск](#запуск)
        - [Настройка](#настройка)
        - [Установка всего необходимого на Windows](#установка-всего-необходимого-на-windows)
        - [Установка всего необходимого на Mac OS](#установка-всего-необходимого-на-mac-os)
        - [Установка make на Linux](#установка-make-на-linux)
    - [Описание игрового процесса](#описание-игрового-процесса)
- [English](#english)
    - [Brief Description](#brief-description)
    - [Setup and Running](#setup-and-running)
        - [Configuration](#configuration)
        - [Windows Installation](#windows-installation)
        - [Mac OS Installation](#mac-os-installation)
        - [Linux Installation](#linux-installation)
    - [Gameplay Description](#gameplay-description)

## Русский

### Краткое описание

Это телеграм бот помощник для квизов в стиле "Своя игра". Он по сути является пультом управления для ведущего и игроков для проведения игр "вживую", находясь в одном помещении. Это решает проблему с определением первого отвечающего игрока, таймером для раунда/ответа и ведениме счёта.

Лично я для игр использую презентации, созданные на основе руководства из этого видео: [https://www.youtube.com/watch?v=xhacW9aT9hA](https://www.youtube.com/watch?v=xhacW9aT9hA)

### Запуск

#### Настройка

Чтобы запустить бота, необходимо получить токен в боте @BotFather

Затем в корне репозитория создать файл .env, его пример можно увидеть в .env.dist:
```text
ROUND_TIMEOUT=15s
ANSWER_TIMEOUT=10s
ROOM_TTL=30m
SESSION_TTL=30m
TELEGRAM_BOT_TOKEN=1234567890:AA...ZZ
LOG_LEVEL=info
```
Здесь в поле TELEGRAM_BOT_TOKEN= прописываем полученный токен.

**ROUND_TIMEOUT -** это время раунда, по истечению которого он автоматически завершается.

**ANSWER_TIMEOUT -** это время ответа, по истечению которого ответ автоматически помечается как **неверный**, у отвечающего игрока отнимаются очки.

**ROOM_TTL -** время жизни комнаты без обновлений. Например, у нас стоит 30m, в таком случае если в комнате ничего не происходило 30 минут, она автоматически удаляется.

**SESSION_TTL -** время жизни сессии пользователя без обновлений. Например, у нас стоит 30m, в таком случае если пользователь ничего не делал 30 минут, его сессия автоматически обнуляется.

**LOG_LEVEL -** это уровень логирования, по дефолту стоит **info**, для отладки можно поставить **debug**. Подробнее про уровни логирования: [https://github.com/rs/zerolog?tab=readme-ov-file#leveled-logging](https://github.com/rs/zerolog?tab=readme-ov-file#leveled-logging)

**TELEGRAM_BOT_TOKEN -** телеграм токен, прописываем сюда полученный токен из @BotFather

Затем открываем Dockerfile и смотрим, чтобы на 13-ой строчке архитектура процессора устройства, на котором будет запущено бот, совпадала с вашей.

Например, вот популярные значения:

Подробнее про то, как задать переменные окружения правильно можно посмотреть здесь: [https://go.dev/doc/install/source#environment](https://go.dev/doc/install/source#environment)

Если собираетесь хостить на системе с обычным 64 битным процессором (amd64, x86-64), то здесь ничего менять не надо.

#### Установка всего необходимого на Windows

Для запуска на Windows потребуется установить Chocolatey, инструкция: [https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe](https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe)

После установки Chocolatey перезапускаем Powershell/cmd

Затем в PowerShell/cmd запускаем команду:
```text
choco install make
```
После этого проверяем успех установки:
```text
make -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
GNU Make 4.4.1
Built for Windows32
Copyright (C) 1988-2023 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later https://gnu.org/licenses/gpl.html
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```
Качаем, устанавливаем и запускаем Docker desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

#### Установка всего необходимого на Mac OS

Устанавливаем Homebrew: [https://brew.sh/](https://brew.sh/)

Там в разделе **Install Homebrew** просто копируем команду и исполняем в терминале.

После этого проверяем успех установки:
```text
brew -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
Homebrew 4.4.8
```
Затем выполняем команду:
```text
brew install make
```
После этого проверяем успех установки:
```text
make -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
GNU Make 3.81
Copyright (C) 2006 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.
There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.
```
Качаем, устанавливаем и запускаем Docker desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

#### Установка make на Linux

Серьёзно? Если вы работаете на Линуксе, то и без меня знаете, что нужно делать.

### Запуск

Выполняем команду в терминале/powershell/cmd:
```text
make build
```
Если сделали всё правильно, ошибок быть не должно. Если есть ошибка, то придётся разобраться с причинами :)

Для запуска необходимо выполнить команду:
```text
make run
```
Для остановки и удаления контейнера можно использовать команду:
```text
make stop
```
Для просмотра логов на своём устройстве рекомендую использовать Docker desktop, для сервера можно ипользовать команду:
```text
make logs
```
Если будете дорабатывать код и пересоздавать образ, то быстро удалить старый можно через команду:
```text
make clean
```
### Описание игрового процесса

Здесь есть 2 роли: админ (ведущий) и игрок.

Админ создаёт комнату, получает её id, отправляет всем игрокам, затем запускает первый раунд.

Если он нажимает "Начать раунд", то комната будет иметь статус "в игре" и к ней уже будет не присоединиться.

При создании раунда необходимо будет ввести количество очков, которые игроки получат за правильный ответ или потеряют при не правильном ответе.

Далее админ должен зачитать/показать вопрос.

При нажатии кнопки "Начать раунд" начинается обратный отсчёт до конца раунда. Игроки получают возможность нажать на кнопку "Ответить".

Если время на раунд заканчивается, то раунд завершается, показывается итоговый счёт, идёт переход к созданию нового раунда.

Если какой-то игрок нажмёт "ответить", то идёт обратный отсчёт до конца ответа, таймер для раунда встаёт на паузу, другие игроки не могут нажать на кнопку "Ответить".

Если время на ответ заканчивается, то ответ автоматически принимается как **неверный,** возвращается таймер раунда с тем временем, на котором он остановился, другие игроки, кроме ответившего, имеют возможность нажать на кнопку "Ответить".

Админ во время ответа игрока получает возможность нажать "Да" или "Нет", чтобы засчитать овет игрока как **верный** и **неверный** соответственно.

При нажатии кнопки "Да" игрок получает плюс то количество очков, которое админ указал при создании раунда, раунд автоматически завершается, показывается итоговый счёт, идёт переход к созданию нового раунда.

При нажатии кнопки "Нет" игрок теряет то количество очков, которое админ указал при создании раунда, возвращается таймер раунда с тем временем, на котором он остановился, другие игроки, кроме ответившего, имеют возможность нажать на кнопку "Ответить".

Счёт игрока может быть отрицательным.

При создании нового раунда админ может нажать на кнопку "Завершить игру".

Если админ нажимает "Завершить игру**",** то показывается итоговый счёт, определяется и показывается имя победителя с максимальным количеством очков. Комната удаляется. Всех перебрасывает на главное меню.

Может быть несколько победителей.

Игрок во время ожидания раунда и во время ожидания ответа другого игрока имеют возможность выйти из комнаты.

В случае, если игрок выходит из комнаты во время игры, которая уже идёт, он уже не сможет присоедениться к ней.

Если игрок выходит из комнаты, пока ещё не началась игра, то он может заново туда зайти.

## English

### Brief Description

This is a Telegram bot assistant for "Jeopardy!"-style quizzes. It essentially serves as a control panel for the host and players to conduct live games in the same physical space. It solves problems like determining the first responder, managing round/answer timers, and keeping score.

Personally, I use presentations created based on the guide from this video: [https://www.youtube.com/watch?v=xhacW9aT9hA](https://www.youtube.com/watch?v=xhacW9aT9hA)

### Setup and Running

#### Configuration

To launch the bot, you need to get a token from @BotFather.

Then create a .env file in the root of the repository (you can see an example in .env.dist):
```text
ROUND_TIMEOUT=15s
ANSWER_TIMEOUT=10s
ROOM_TTL=30m
SESSION_TTL=30m
TELEGRAM_BOT_TOKEN=1234567890:AA...ZZ
LOG_LEVEL=info
```
Enter the obtained token in the TELEGRAM_BOT_TOKEN field.

**ROUND_TIMEOUT -** round duration, after which the round automatically ends.

**ANSWER_TIMEOUT -** answer time limit, after which the answer is automatically marked as **incorrect** and points are deducted from the responding player.

**ROOM_TTL -** room lifetime without updates. For example, if set to 30m, the room will be automatically deleted after 30 minutes of inactivity.

**SESSION_TTL -** user session lifetime without updates. For example, if set to 30m, the user session will be automatically reset after 30 minutes of inactivity.

**LOG_LEVEL -** logging level, default is **info**, for debugging you can set **debug**. More about logging levels: [https://github.com/rs/zerolog?tab=readme-ov-file#leveled-logging](https://github.com/rs/zerolog?tab=readme-ov-file#leveled-logging)

**TELEGRAM_BOT_TOKEN -** Telegram token, enter the token obtained from @BotFather here.

Then open Dockerfile and check that the processor architecture on line 13 matches your device's architecture.

For example, here are common values:

More about setting environment variables correctly: [https://go.dev/doc/install/source#environment](https://go.dev/doc/install/source#environment)

If you're hosting on a system with a regular 64-bit processor (amd64, x86-64), you don't need to change anything here.

#### Windows Installation

For Windows, you'll need to install Chocolatey: [https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe](https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe)

After installing Chocolatey, restart PowerShell/cmd.

Then run in PowerShell/cmd:
```text
choco install make
```
After installation, verify success:
```text
make -v
```
If you see something like this, everything is OK:
```text
GNU Make 4.4.1
Built for Windows32
Copyright (C) 1988-2023 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later https://gnu.org/licenses/gpl.html
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```
Download, install and run Docker Desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

#### Mac OS Installation

Install Homebrew: [https://brew.sh/](https://brew.sh/)

In the **Install Homebrew** section, just copy the command and execute it in terminal.

After installation, verify success:
```text
brew -v
```
If you see something like this, everything is OK:
```text
Homebrew 4.4.8
```
Then run:
```text
brew install make
```
After installation, verify success:
```text
make -v
```
If you see something like this, everything is OK:
```text
GNU Make 3.81
Copyright (C) 2006 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.
There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.
```
Download, install and run Docker Desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

#### Linux Installation

Seriously? If you're using Linux, you know what to do without me.

### Running

Execute in terminal/PowerShell/cmd:
```text
make build
```
If everything was done correctly, there should be no errors. If there's an error, you'll need to figure out the cause :)

To run:
```text
make run
```
To stop and remove the container:
```text
make stop
```
To view logs on your local device, I recommend using Docker Desktop. For servers, you can use:
```text
make logs
```
If you're modifying the code and recreating the image, you can quickly remove the old one with:
```text
make clean
```
### Gameplay Description

There are 2 roles: admin (host) and player.

The admin creates a room, gets its ID, sends it to all players, then starts the first round.

If they click "Start round", the room status changes to "in game" and no one can join anymore.

When creating a round, you need to specify the point value that players will get for a correct answer or lose for an incorrect answer.

Then the admin should read/show the question.

When "Start round" is clicked, the countdown begins. Players get the option to click "Answer".

If the round time runs out, the round ends, the final score is shown, and the game moves to creating a new round.

If a player clicks "Answer", the answer countdown begins, the round timer pauses, and other players can't click "Answer".

If the answer time runs out, the answer is automatically marked as **incorrect**, the round timer resumes from where it paused, and other players (except the one who answered) can click "Answer".

During a player's answer, the admin can click "Yes" or "No" to mark the answer as **correct** or **incorrect** respectively.

Clicking "Yes" gives the player the points specified when creating the round, automatically ends the round, shows the final score, and moves to creating a new round.

Clicking "No" deducts the specified points from the player, resumes the round timer, and allows other players (except the one who answered) to click "Answer".

A player's score can be negative.

When creating a new round, the admin can click "End game".

Clicking "End game" shows the final score, determines and shows the winner(s) with the highest score, deletes the room, and returns everyone to the main menu.

There can be multiple winners.

Players can leave the room while waiting for a round or another player's answer.

If a player leaves during an ongoing game, they can't rejoin.

If a player leaves before the game starts, they can rejoin.