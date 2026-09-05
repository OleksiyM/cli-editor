import random

def hello_world_or_name(name: str | None = None) -> str:
    return f"Hello, {name if name is not None else 'World'}!"

def read_name() -> str | None:
    name = input("Enter your name: ").strip().capitalize()
    return name or None

def generate_random_number(min_value: int, max_value: int) -> int:
    return random.randint(min_value, max_value)

def select_random_chat(chats: list[str], chat_counts: dict) -> str:
    if not chats:
        return "No chats available."
    choice = random.choice(chats)
    chat_counts[choice] = chat_counts.get(choice, 0) + 1  # Считаем, сколько раз выпал
    duration = generate_random_number(1, 10)  # Длительность чата в минутах
    return f"{choice} for {duration} minutes"

if __name__ == "__main__":
    hello_name = read_name()
    print(hello_world_or_name(hello_name))
    chats = ["Python", "Data Science", "English", "Brain Storm", "Just Chatting"]
    chat_counts = {}  # Словарь для подсчета
    for _ in range(3):  # Три "раунда" чата, как три куплета
        print(f"Now selected: {select_random_chat(chats, chat_counts)}")
    print("\nChat stats:", chat_counts)