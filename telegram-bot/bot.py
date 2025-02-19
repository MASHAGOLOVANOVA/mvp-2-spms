# -*- coding: utf-8 -*-
"""Этот модуль реализует функциональность бота Telegram."""

import json
from datetime import datetime, timedelta
import telebot
import requests
from requests.exceptions import RequestException
from telegram.constants import ParseMode

HOST_URL = "http://localhost:8080"

CLIENT_URL = "http://localhost:3000"

class SessionManager:
    """Класс для получении инфы о сессии"""

    def __init__(self):
        """Класс для получении инфы о сессии"""
        self.session_token = None

    def set_session_token(self, token):
        """функция для апдейта токена"""
        self.session_token = token

    def get_headers(self):
        """функция для получения хедеров запросов"""
        return {
            "Content-Type": "application/json",
            "Bot-Token": "7772483926:AAFkT_nibrVHwZmlJajxbXRU4Wxe_b7t_RI",
            "tuna-skip-browser-warning": "please",
            "Session-Id": self.session_token,
        }


sessionManager = SessionManager()


def update_session_token(new_token):
    """# Функция для обновления session_token"""
    sessionManager.set_session_token(new_token)


BOT_TOKEN = "7772483926:AAFkT_nibrVHwZmlJajxbXRU4Wxe_b7t_RI"
# Создаем экземпляр бота
bot = telebot.TeleBot(BOT_TOKEN)


@bot.message_handler(commands=["start"])
def start_message(message):
    """Хендлер команды start"""
    # Создаем кнопку для запроса номера телефона
    keyboard = telebot.types.ReplyKeyboardMarkup(one_time_keyboard=True)
    button = telebot.types.KeyboardButton("Отправить номер телефона", request_contact=True)
    keyboard.add(button)

    bot.send_message(
        message.chat.id,
        """Привет! \n
Это бот Системы для управления студенческими проектами!
Пожалуйста, отправьте свой номер телефона.""",
        reply_markup=keyboard,
    )


@bot.message_handler(content_types=["contact"])
def handle_contact(message):
    """Хендлер contact"""
    contact = message.contact
    phone_number = contact.phone_number  # Получаем номер телефона
    bot.send_message(message.chat.id, f"Спасибо! Ваш номер телефона: {phone_number}")

    # Создаем данные для отправки на сервер
    credentials = {"phone_number": phone_number}
    verify_number(message, credentials)


def show_main_menu(chat_id):
    """функция открытия главного меню"""
    # Создаем главное меню
    keyboard = telebot.types.ReplyKeyboardMarkup(resize_keyboard=True)
    button_projects = telebot.types.KeyboardButton("Мои проекты")
    button_meetings = telebot.types.KeyboardButton("Мои встречи")
    button_add_project = telebot.types.KeyboardButton("Добавить проект")  # Новая кнопка

    has_planner = get_google_planner()
    if has_planner is not None:
        keyboard.add(button_projects, button_meetings, button_add_project)
    else:
        bot.send_message(
            chat_id,
            f"""К сожалению Вам недоступно расписание встреч!\n
Чтобы пользоваться расписанием, подключите Google Calendar из веб-приложения:
\n{CLIENT_URL}/profile/integrations""",
        )
        keyboard.add(button_projects, button_add_project)

    bot.send_message(chat_id, "Выберите действие:", reply_markup=keyboard)


@bot.message_handler(func=lambda message: message.text == "Мои проекты")
def handle_projects(message):
    """Хендлер проектов"""
    try:
        response = requests.get(
            f"{HOST_URL}/api/v1/projects/",
            headers=sessionManager.get_headers(),
            timeout=10,
        )

        # Проверяем статус ответа
        if response.status_code == 200:
            response_data = response.json()  # Предполагаем, что ответ в формате JSON

            projects = response_data.get(
                "projects", []
            )  # Предполагаем, что проекты находятся в ключе 'projects'
            if projects:  # Проверяем, есть ли проекты в списке
                for project in projects:
                    # Создаем карточку с кнопкой для каждого проекта
                    project_card = (
                        f"""Тема: {project['theme']}\nГод: {project['year']}\n"""
                    )

                    # Создаем кнопку для каждого проекта
                    markup = telebot.types.InlineKeyboardMarkup()
                    button = telebot.types.InlineKeyboardButton(
                        "Подробнее", callback_data=f"project_{project['id']}"
                    )
                    markup.add(button)

                    # Отправляем сообщение с карточкой и кнопкой
                    bot.send_message(message.chat.id, project_card, reply_markup=markup)
            else:
                bot.send_message(message.chat.id, "У вас нет проектов.")
        else:
            bot.send_message(
                message.chat.id,
                f"Ошибка при получении проектов: {response.status_code}",
            )
    except RequestException as e:
        bot.send_message(message.chat.id, f"Ошибка: {str(e)}")


@bot.message_handler(func=lambda message: message.text == "Мои встречи")
def handle_meetings(message):
    """Хендлер встреч"""
    meetings = get_schedule()
    if len(meetings) > 0:
        response = format_meetings(group_meetings_by_day(meetings))
        for d in response:
            bot.send_message(message.chat.id, d, parse_mode=ParseMode.MARKDOWN)
    else:
        bot.send_message(message.chat.id, "Встречи не назначены")


days_translation = {
    "Monday": "Понедельник",
    "Tuesday": "Вторник",
    "Wednesday": "Среда",
    "Thursday": "Четверг",
    "Friday": "Пятница",
    "Saturday": "Суббота",
    "Sunday": "Воскресенье",
}


def group_meetings_by_day(meetings):
    """функция для группировки встреч"""
    grouped = {}
    for meeting in meetings:
        meeting_time = datetime.fromisoformat(meeting["time"].replace("Z", "+00:00"))
        day = days_translation.get(meeting_time.strftime("%A"))  # Получаем день недели
        date = meeting_time.strftime("%d.%m.%Y")
        day += f", {date}"
        if day not in grouped:
            grouped[day] = []
        grouped[day].append(meeting)
    return grouped


def format_meetings(grouped_meetings):
    """функция для форматирования встреч"""
    alldays = []
    for day, meetings in grouped_meetings.items():
        response = f"*{day}*\n\n"  # Заголовок дня недели
        for meeting in meetings:
            start_time = datetime.fromisoformat(meeting["time"].replace("Z", "+00:00"))
            end_time = start_time + timedelta(
                hours=1
            )  # Добавляем 1 час к начальному времени

            # Форматируем время
            formatted_start_time = start_time.strftime("%H:%M")
            formatted_end_time = end_time.strftime("%H:%M")
            response += f"{formatted_start_time}"
            response += f" - {formatted_end_time}\nНазвание: {meeting['name']}\n"
            response += f"Описание: {meeting['description']}\n"
            response += f"Студент: {meeting['student']['name']},"
            response += f" Курс: {meeting['student']['cource']}\n"
            response += f"{'Онлайн' if meeting['is_online'] else 'Оффлайн'}\n\n"
        response += "\n"
        alldays.append(response)
    return alldays


def get_schedule():
    """функция для получения расписания"""
    current_time = datetime.utcnow()
    start_of_day = current_time.replace(hour=0, minute=0, second=0, microsecond=0)

    iso_format_time = start_of_day.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"

    # Формируем URL с параметром from
    url = f"{HOST_URL}/api/v1/meetings?from={iso_format_time}"

    # Выполняем GET-запрос
    response = requests.get(url, headers=sessionManager.get_headers(), timeout=10)
    if response.status_code == 200:
        response_data = response.json()
        print(response_data)
        meetings = response_data.get("meetings", [])
        return meetings
    return []


def verify_number(message, credentials):
    """функция для поиска номера в системе"""
    try:
        bot.send_message(message.chat.id, "Проверяем регистрацию...")
        response = requests.post(
            HOST_URL + "/api/v1/auth/bot/signinuser",
            json=credentials,
            headers=sessionManager.get_headers(),
            timeout=10,
        )
        if response.status_code == 200:
            bot.send_message(message.chat.id, "Мы Вас нашли!")

            response_data = json.loads(response.text)  # Используем json.loads()

            # Обновляем session_token, если он присутствует в ответе
            if "session_token" in response_data:
                update_session_token(response_data["session_token"])
                professor = get_account()
                bot.send_message(message.chat.id, f'Здравствуйте, {professor["name"]}!')
                cloud_drive = get_cloud_drive()
                if cloud_drive is not None:
                    show_main_menu(message.chat.id)
                else:
                    bot.send_message(
                        message.chat.id,
                        f"""Чтобы воспользоваться функциями бота подключите
 Google Drive из веб-приложения:\n{CLIENT_URL}/profile/integrations""",
                    )
            else:
                print("session_token не найден в ответе")
        else:
            bot.send_message(
                message.chat.id,
                "Произошла ошибка при поиске пользователя по номеру телефона.",
            )

    except RequestException as e:
        bot.send_message(message.chat.id, f"Ошибка сети: {str(e)}")


def get_account():
    """функция для получения аккаунта"""
    response = requests.get(
        f"{HOST_URL}/api/v1/account", headers=sessionManager.get_headers(), timeout=10
    )
    if response.status_code == 200:
        account = response.json()
        return account
    return []


def get_integrations():
    """функция для получения интеграций"""
    integrations_response = requests.get(
        f"{HOST_URL}/api/v1/account/integrations",
        headers=sessionManager.get_headers(),
        timeout=10,
    )
    return integrations_response


def get_cloud_drive():
    """функция для получения интеграции с диском"""
    integrations = get_integrations()
    try:
        response_json = integrations.json()
        if "cloud_drive" in response_json:
            return response_json["cloud_drive"]
        print("Cloud Drive не найден.")
        return None
    except ValueError as e:
        print(f"Ошибка при обработке JSON: {e}")
        return None


def get_google_planner():
    """функция для получения интеграции с планнером"""
    integrations = get_integrations()
    print(integrations.text)
    try:
        response_json = integrations.json()
        # Проверяем наличие "cloud_drive"
        if "planner" in response_json:
            if "planner_name" in response_json["planner"]:
                return response_json["planner"]["planner_name"]
            return None
        print("Cloud Calendar не найден.")
        return None
    except ValueError as e:
        print(f"Ошибка при обработке JSON: {e}")
        return None


@bot.message_handler(func=lambda message: message.text == "Добавить проект")
def add_project(message):
    """Хендлер команды Добавить проект"""
    students = get_students()

    if not students:
        bot.send_message(
            message.chat.id,
            """Нет доступных студентов для выбора.
Вы можете добавить нового студента. Введите имя нового студента:""",
        )
        bot.register_next_step_handler(message, add_student_name)
        return
    # Создаем клавиатуру для выбора студента
    keyboard = telebot.types.ReplyKeyboardMarkup(resize_keyboard=True)
    for student in students:
        keyboard.add(
            telebot.types.KeyboardButton(
                student["surname"] + " " + student["name"] + " " + student["middlename"]
            )
        )  # Предполагаем, что у студента есть поле 'name'

    bot.send_message(
        message.chat.id, "Выберите студента для проекта:", reply_markup=keyboard
    )
    bot.register_next_step_handler(message, process_student_selection)


def process_student_selection(message):
    """функция для выбора студента"""
    student_name = message.text
    # Здесь вы можете добавить проверку на наличие выбранного студента в списке
    # Например, если у вас есть список студентов в виде словаря
    students = get_students()

    selected_student = next(
        (s for s in students if s["name"] == student_name.split(" ")[0]), None
    )

    if selected_student is None:
        bot.send_message(
            message.chat.id,
            "Выбранный студент не найден. Пожалуйста, попробуйте снова.",
        )
        show_main_menu(message.chat.id)
        return

    bot.send_message(message.chat.id, "Введите тему нового проекта:")
    bot.register_next_step_handler(
        message, lambda msg: process_project_theme(msg, selected_student)
    )


def process_project_theme(message, student):
    """функция для выбора темы проекта"""
    project_theme = message.text
    bot.send_message(message.chat.id, "Введите год проекта (число):")
    bot.register_next_step_handler(
        message, lambda msg: process_project_year(msg, student, project_theme)
    )


def process_project_year(message, student, project_theme):
    """функция для выбора года проекта"""
    try:
        project_year = int(message.text)
        bot.send_message(message.chat.id, "Введите владельца репозитория (логин):")
        bot.register_next_step_handler(
            message,
            lambda msg: process_repo_owner(msg, student, project_theme, project_year),
        )
    except ValueError:
        bot.send_message(
            message.chat.id, "Год должен быть числом. Пожалуйста, попробуйте снова."
        )
        process_project_year(message, student, project_theme)


def process_repo_owner(message, student, project_theme, project_year):
    """функция для выбора владельца проекта"""
    repo_owner = message.text
    bot.send_message(message.chat.id, "Введите имя репозитория:")
    bot.register_next_step_handler(
        message,
        lambda msg: process_repository_name(
            msg, student, project_theme, project_year, repo_owner
        ),
    )


def process_repository_name(message, student, project_theme, project_year, repo_owner):
    """функция для ввода имени репозитория"""
    repository_name = message.text

    response = requests.post(
        f"{HOST_URL}/api/v1/projects/add",
        json={
            "theme": project_theme,
            "student_id": student["id"],
            "year": project_year,
            "repository_owner_login": repo_owner,
            "repository_name": repository_name,
        },
        headers=sessionManager.get_headers(),
        timeout=10,
    )

    if response.status_code == 200:
        bot.send_message(
            message.chat.id,
            f'Проект "{project_theme}" успешно добавлен для студента "{student["name"]}"!',
        )
    else:
        bot.send_message(
            message.chat.id, f"Ошибка при добавлении проекта: {response.status_code}"
        )

    # Вернуться в главное меню после добавления проекта
    show_main_menu(message.chat.id)


# Предполагается, что у вас есть функция для получения списка студентов
def get_students():
    """функция для получения студентов"""
    response = requests.get(
        f"{HOST_URL}/api/v1/students", headers=sessionManager.get_headers(), timeout=10
    )
    if response.status_code == 200:
        response_data = response.json()
        students = response_data.get("students", [])
        return students  # Возвращает список студентов в формате JSON
    return []


def get_educational_programmes():
    """функция для получения учебных программ"""
    response = requests.get(
        f"{HOST_URL}/api/v1/universities/1/edprogrammes/",
        headers=sessionManager.get_headers(),
        timeout=10,
    )
    if response.status_code == 200:
        response_data = response.json()
        educational_programmes = response_data.get("programmes", [])
        return educational_programmes  # Возвращает список образовательных программ в формате JSON
    return []


def add_student_name(message):
    """функция для ввода имени студента"""
    student_name = message.text
    bot.send_message(message.chat.id, "Введите фамилию нового студента:")
    bot.register_next_step_handler(
        message, lambda msg: add_student_surname(msg, student_name)
    )


def add_student_surname(message, student_name):
    """функция для ввода фамилии студента"""
    student_surname = message.text
    bot.send_message(message.chat.id, "Введите отчество нового студента:")
    bot.register_next_step_handler(
        message, lambda msg: add_student_middlename(msg, student_name, student_surname)
    )


def add_student_middlename(message, student_name, student_surname):
    """функция для ввода отчества студента"""
    student_middlename = message.text
    bot.send_message(message.chat.id, "Введите курс нового студента (число):")
    bot.register_next_step_handler(
        message,
        lambda msg: add_student_course(
            msg, student_name, student_surname, student_middlename
        ),
    )


def add_student_course(message, student_name, student_surname, student_middlename):
    """функция для ввода курса студента"""
    try:
        student_course = int(message.text)
        educational_programmes = get_educational_programmes()
        if not educational_programmes:
            bot.send_message(
                message.chat.id,
                "Нет доступных образовательных программ. Попробуйте позже.",
            )
            show_main_menu(message.chat.id)
            return

        # Создаем клавиатуру для выбора образовательной программы
        keyboard = telebot.types.ReplyKeyboardMarkup(resize_keyboard=True)
        for programme in educational_programmes:
            keyboard.add(telebot.types.KeyboardButton(programme["name"]))

        bot.send_message(
            message.chat.id,
            "Выберите образовательную программу:",
            reply_markup=keyboard,
        )
        bot.register_next_step_handler(
            message,
            lambda msg: add_student_programme(
                msg, student_name, student_surname, student_middlename, student_course
            ),
        )
    except ValueError:
        bot.send_message(
            message.chat.id, "Курс должен быть числом. Пожалуйста, попробуйте снова."
        )
        add_student_course(message, student_name, student_surname, student_middlename)


def add_student_programme(
    message, student_name, student_surname, student_middlename, student_course
):
    """функция для ввода уч программы студента"""
    selected_programme_name = message.text
    educational_programmes = get_educational_programmes()

    selected_programme = next(
        (ep for ep in educational_programmes if ep["name"] == selected_programme_name),
        None,
    )
    if selected_programme is None:
        bot.send_message(
            message.chat.id,
            "Выбранная программа не найдена. Пожалуйста, попробуйте снова.",
        )
        show_main_menu(message.chat.id)
        return

    # Создаем нового студента
    new_student_data = {
        "name": student_name,
        "surname": student_surname,
        "middlename": student_middlename,
        "cource": student_course,
        "education_programme_id": selected_programme["id"],
    }
    response = requests.post(
        f"{HOST_URL}/api/v1/students/add",
        json=new_student_data,
        headers=sessionManager.get_headers(),
        timeout=10,
    )
    if response.status_code == 200:
        bot.send_message(
            message.chat.id,
            f'Студент "{student_name} {student_surname}" успешно добавлен!',
        )
        add_project(
            message
        )  # После добавления студента можно снова вызвать функцию добавления проекта
    else:
        bot.send_message(
            message.chat.id,
            f"Ошибка при добавлении студента: {response.status_code}. Попробуйте снова.",
        )
        show_main_menu(message.chat.id)


@bot.callback_query_handler(func=lambda call: call.data.startswith("project_"))
def handle_project_details(call):
    """Хендлер для просмотра информации о проекте"""
    project_id = call.data.split("_")[1]
    headers = sessionManager.get_headers()
    response = requests.get(
        f"{HOST_URL}/api/v1/projects/{project_id}", headers=headers, timeout=10
    )

    if response.status_code == 200:
        project_details = response.json()

        student = project_details["student"]
        print(student)
        student_str = f"{student['surname']} {student['name']} {student['middlename']}"
        theme = project_details["theme"]
        details_message = (
            "*Тема:* "
            + theme
            + "\n"
            + "*Год:* "
            + str(project_details["year"])
            + "\n"
            + "*Студент:* "
            + student_str
            + "\n"
            + "*Статус проекта:* "
            + project_details["status"]
            + "\n"
            + "*Стадия работы:* "
            + project_details["stage"]
            + "\n"
            + "*Ссылка на Google Drive:* [Перейти к папке]("
            + project_details["cloud_folder_link"]
            + ")\n"
        )
        # Создание кнопок
        markup = telebot.types.InlineKeyboardMarkup()
        button1 = telebot.types.InlineKeyboardButton(
            "Статистика", callback_data=f"statistics_project_{project_details['id']}"
        )
        button2 = telebot.types.InlineKeyboardButton(
            "Коммиты", callback_data=f"commits_project_{project_details['id']}"
        )
        button3 = telebot.types.InlineKeyboardButton(
            "Задания", callback_data=f"tasks_project_{project_details['id']}"
        )
        button4 = telebot.types.InlineKeyboardButton(
            "Назначить задание",
            callback_data=f"add_task_project_{project_details['id']}",
        )
        button5 = telebot.types.InlineKeyboardButton(
            "Назначить встречу",
            callback_data="add_meeting_project_"
            + f"{project_details['id']}_student_{project_details['student']['id']}",
        )
        if get_repohub() is not None:
            markup.add(button1, button2, button3, button4, button5)
        else:
            bot.send_message(
                call.message.chat.id,
                f"""Вам недоступны коммиты проекта, подключите интеграцию
с Github в личном кабинете в веб-приложении: <a href='{CLIENT_URL}/profile/integrations'>
Перейти к интеграциям</a>""",
                parse_mode="HTML",
            )
            markup.add(button1, button3, button4, button5)
        bot.send_message(
            call.message.chat.id,
            details_message,
            reply_markup=markup,
            parse_mode="Markdown",
        )
    else:
        bot.send_message(
            call.message.chat.id,
            f"Ошибка при получении деталей проекта: {response.status_code}",
        )
    bot.edit_message_reply_markup(call.message.chat.id, call.message.message_id)


def format_statistics_message(statistics):
    """Форматирует сообщение со статистикой проекта."""
    total_meetings = statistics.get("total_meetings", 0)
    total_tasks = statistics.get("total_tasks", 0)
    tasks_done = statistics.get("tasks_done", 0)
    tasks_done_percent = statistics.get("tasks_done_percent", 0)

    stats_message = (
        "*Статистика по проекту:*\n\n"
        f"📅 *Общее количество встреч:* {total_meetings}\n"
        f"📋 *Общее количество задач:* {total_tasks}\n"
        f"✅ *Завершенные задачи:* {tasks_done} ({tasks_done_percent}%)\n\n"
    )

    grades = statistics.get("grades", {})
    if grades:
        stats_message += format_grades(grades)
    else:
        stats_message += "Оценки отсутствуют.\n"
    return stats_message

def format_grades(grades):
    """Форматирует сообщение с оценками."""
    defence_grade = grades.get("defence_grade", "Нет оценки")
    supervisor_grade = grades.get("supervisor_grade", "Нет оценки")
    final_grade = grades.get("final_grade", "Нет оценки")
    supervisor_review = grades.get("supervisor_review", {})
    grades_message = (
        "*Оценки:*\n"
        f"🎓 *Защита:* {defence_grade}\n"
        f"👨‍🏫 *Оценка руководителя:* {supervisor_grade}\n"
        f"🏆 *Итоговая оценка:* {final_grade}\n\n"
    )
    if supervisor_review:
        review_criterias = supervisor_review.get("criterias", [])
        if review_criterias:
            grades_message += "*Критерии оценки:*\n"
            for criteria in review_criterias:
                criteria_name = criteria.get("criteria", "Не указано")
                criteria_grade = criteria.get("grade", "Не указано")
                criteria_weight = criteria.get("weight", "Не указано")
                grades_message += f"- {criteria_name}: Оценка {criteria_grade} (Вес: {criteria_weight})\n"

    return grades_message


@bot.callback_query_handler(
    func=lambda call: call.data.startswith("statistics_project_")
)
def handle_project_statisctics(call):
    """Хендлер для просмотра статистики по проекту"""
    project_id = call.data.split("_")[2]
    response = requests.get(
        f"{HOST_URL}/api/v1/projects/{project_id}/statistics",
        headers=sessionManager.get_headers(),
        timeout=10,
    )

    if response.status_code == 200:
        statistics = response.json()
        stats_message = format_statistics_message(statistics)
        bot.send_message(call.message.chat.id, stats_message, parse_mode="Markdown")
    else:
        bot.send_message(
            call.message.chat.id,
            f"Ошибка при получении статистики проекта: {response.status_code}",
        )

def format_commit_message(commit):
    """Форматирует сообщение о коммите."""
    commit_sha = commit.get("commit_sha", "Не указано")
    message = commit.get("message", "Не указано")
    date_created = commit.get("date_created", "Не указано")
    created_by = commit.get("created_by", "Не указано")
    formatted_date = datetime.fromisoformat(date_created[:-1]).strftime("%Y-%m-%d %H:%M:%S")
    return (
        f"🔹 *SHA:* {commit_sha}\n"
        f"📝 *Сообщение:* {message}\n"
        f"📅 *Дата создания:* {formatted_date}\n"
        f"👤 *Создано пользователем:* {created_by}\n\n"
    )

@bot.callback_query_handler(func=lambda call: call.data.startswith("commits_project_"))
def handle_project_commits(call):
    """Хендлер для просмотра коммитов по проекту"""
    project_id = call.data.split("_")[2]
    current_time = datetime.utcnow() - timedelta(days=30)
    month_ago = current_time.replace(hour=0, minute=0, second=0, microsecond=0)

    iso_format_time = month_ago.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
    url = f"{HOST_URL}/api/v1/projects/{project_id}/commits?from={iso_format_time}"
    response = requests.get(url, headers=sessionManager.get_headers(), timeout=10)
    if response.status_code == 200:
        commits_data = response.json()
        commits = commits_data.get("commits", [])

        if commits:
            commits_message = "*Коммиты проекта:*\n\n"
            for commit in commits:
                commits_message += format_commit_message(commit)

            bot.send_message(call.message.chat.id, commits_message, parse_mode="Markdown")
        else:
            bot.send_message(call.message.chat.id, "Коммиты отсутствуют.")
    else:
        bot.send_message(
            call.message.chat.id,
            f"Ошибка при получении коммитов проекта: {response.status_code}",
        )


def get_repohub():
    """функция для получения интеграции с гитхабом"""
    integrations = get_integrations()
    if integrations.status_code == 200:
        integrations_data = integrations.json()  # Преобразуем в JSON
        if len(integrations_data["repo_hubs"]) > 0:
            return integrations_data["repo_hubs"]
    print(f"Ошибка при получении интеграций: {integrations.status_code}")
    return None


@bot.callback_query_handler(func=lambda call: call.data.startswith("add_task_project_"))
def handle_project_new_task(call):
    """Хендлер для добавления задачи"""
    project_id = call.data.split("_")[3]

    bot.send_message(call.message.chat.id, "Введите название задачи:")
    bot.register_next_step_handler(call.message, process_task_name, project_id)


def process_task_name(message, project_id):
    """функция для получения названия задачи"""
    task_name = message.text  # Получаем название задачи

    bot.send_message(message.chat.id, "Введите описание задачи:")
    bot.register_next_step_handler(
        message, process_task_description, project_id, task_name
    )


def process_task_description(message, project_id, task_name):
    """функция для получения описания задачи"""
    task_description = message.text  # Получаем описание задачи

    bot.send_message(
        message.chat.id, "Введите дедлайн задачи (в формате dd.mm.YYYY HH:MM):"
    )
    bot.register_next_step_handler(
        message, process_task_deadline, project_id, task_name, task_description
    )


def process_task_deadline(message, project_id, task_name, task_description):
    """функция для получения дедлайна задачи"""
    task_deadline_input = message.text  # Получаем дедлайн задачи
    try:
        deadline_datetime = datetime.strptime(task_deadline_input, "%d.%m.%Y %H:%M")

        new_task_data = {
            "name": task_name,
            "description": task_description,
            "deadline": deadline_datetime.isoformat()
            + "Z",  # Преобразуем в строку ISO 8601
        }

        url = f"{HOST_URL}/api/v1/projects/{project_id}/tasks/add"

        response = requests.post(
            url, json=new_task_data, headers=sessionManager.get_headers(), timeout=10
        )

        if response.status_code == 200:
            bot.send_message(message.chat.id, "Задача успешно добавлена!")
        else:

            bot.send_message(
                message.chat.id, f"Ошибка при добавлении задачи: {response.text}"
            )

    except ValueError:
        bot.send_message(
            message.chat.id,
            "Неверный формат даты. Пожалуйста, используйте формат dd.mm.YYYY HH:MM.",
        )


@bot.callback_query_handler(func=lambda call: call.data.startswith("tasks_project_"))
def handle_project_tasks(call):
    """Хендлер для получения задач проекта"""
    project_id = call.data.split("_")[2]
    url = f"{HOST_URL}/api/v1/projects/{project_id}/tasks"

    response = requests.get(url, headers=sessionManager.get_headers(), timeout=10)

    if response.status_code == 200:
        tasks_data = response.json()
        tasks = tasks_data.get("tasks", [])

        if tasks:
            tasks_message = "*Задания по проекту:*\n\n"
            for task in tasks:
                task_id = task.get("id", "Не указано")
                task_name = task.get("name", "Не указано")
                task_description = task.get("description", "Не указано")
                task_deadline = task.get("deadline", "Не указано")
                task_status = task.get("status", "Не указано")
                cloud_folder_link = task.get("cloud_folder_link", "Не указано")

                formatted_deadline = datetime.fromisoformat(
                    task_deadline[:-1]
                ).strftime("%Y-%m-%d %H:%M:%S")

                tasks_message += f"🔹 *ID:* `{task_id}`\n"
                tasks_message += f"📝 *Название:* {task_name}\n"
                tasks_message += f"📜 *Описание:* {task_description}\n"
                tasks_message += f"📅 *Дедлайн:* {formatted_deadline}\n"
                tasks_message += f"🔄 *Статус:* {task_status}\n"
                tasks_message += (
                    f"📂 *Ссылка на папку:* [Google Drive]({cloud_folder_link})\n\n"
                )

            bot.send_message(call.message.chat.id, tasks_message, parse_mode="Markdown")
        else:
            bot.send_message(call.message.chat.id, "Задания отсутствуют.")
    else:
        bot.send_message(
            call.message.chat.id,
            f"Ошибка при получении заданий: {response.status_code}",
        )


@bot.callback_query_handler(
    func=lambda call: call.data.startswith("add_meeting_project_")
)
def handle_project_new_meeting(call):
    """Хендлер для добавления встречи по проекту"""
    project_id = call.data.split("_")[3]
    student_id = call.data.split("_")[5]

    bot.send_message(call.message.chat.id, "Введите название встречи:")
    bot.register_next_step_handler(
        call.message, process_meeting_name, project_id, student_id
    )


def process_meeting_name(message, project_id, student_id):
    """функция для получения названия встречи"""
    name = message.text  # Получаем название

    bot.send_message(message.chat.id, "Введите описание встречи:")
    bot.register_next_step_handler(
        message, process_meeting_description, project_id, student_id, name
    )


def process_meeting_description(message, project_id, student_id, name):
    """функция для получения описания встречи"""
    desc = message.text  # Получаем название

    bot.send_message(message.chat.id, "Введите время встречи:")
    bot.register_next_step_handler(
        message, process_meeting_time, project_id, student_id, name, desc
    )


def process_meeting_time(message, project_id, student_id, name, desc):
    """функция для получения времени встречи"""
    time = message.text  # Получаем название встречи

    bot.send_message(
        message.chat.id,
        "Выберите формат встречи:",
        reply_markup=get_meeting_format_markup(),
    )
    meeting_options = {"name": name, "desc": desc}
    bot.register_next_step_handler(
        message, process_meeting_format, project_id, student_id, meeting_options, time
    )


def get_meeting_format_markup():
    """функция для получения доски для выбора формата встречи"""
    markup = telebot.types.ReplyKeyboardMarkup(one_time_keyboard=True)
    button_online = telebot.types.KeyboardButton("Онлайн")
    button_offline = telebot.types.KeyboardButton("Оффлайн")
    markup.add(button_online, button_offline)
    return markup


def process_meeting_format(message, project_id, student_id, meeting_options, time):
    """функция для получения формата встречи"""
    meeting_format = message.text  # Получаем формат встречи

    if meeting_format not in ["Онлайн", "Оффлайн"]:
        bot.send_message(
            message.chat.id,
            "Пожалуйста, выберите корректный формат встречи: Онлайн или Оффлайн.",
        )
        return  # Завершаем выполнение функции, если формат некорректный
    try:
        online = meeting_format == "Онлайн"  # Устанавливаем значение is_online
        iso_time = (datetime.strptime(time, "%d.%m.%Y %H:%M")).isoformat()
        # Формируем данные для новой встречи
        new_meeting_data = {
            "name": meeting_options["name"],
            "description": meeting_options["desc"],
            "project_id": int(project_id),
            "student_participant_id": int(student_id),
            "is_online": online,
            "meeting_time": iso_time + "Z",  # Преобразуем в строку ISO 8601
        }
        url = f"{HOST_URL}/api/v1/meetings/add"
        response = requests.post(
            url, json=new_meeting_data, headers=sessionManager.get_headers(), timeout=10
        )

        if response.status_code == 200:
            bot.send_message(message.chat.id, "Встреча успешно добавлена!")
        else:
            print(response.text)
            bot.send_message(
                message.chat.id,
                f"Ошибка при добавлении встречи: {response.status_code}",
            )
    except ValueError:
        bot.send_message(
            message.chat.id,
            "Неверный формат даты. Пожалуйста, используйте формат YYYY-MM-DD HH:MM.",
        )

bot.polling(none_stop=True, interval=0)
