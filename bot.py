import asyncio
import logging
import threading
import time

import schedule
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import ApplicationBuilder, ContextTypes, MessageHandler, \
    filters, CallbackQueryHandler, CommandHandler, ConversationHandler

import poller
import storage

BOT_TOKEN = ''
POLLER: poller.Poller
STORAGE: storage.AbstractStorage

VEHICLE_CHOICE, MIN_COUNT_INPUT = range(2)


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text('Hello! I am your bot. Use /getdata to fetch the latest data.')


async def add_monitoring(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    vehicles = [s for s in await POLLER.get_data() if s.current_count > 0]

    keyboard = [[InlineKeyboardButton(v.name, callback_data=v.name)] for v in vehicles]
    reply_markup = InlineKeyboardMarkup(keyboard)

    # Send the list of buttons to the user
    await update.message.reply_text('Choose a vehicle:', reply_markup=reply_markup)

    return VEHICLE_CHOICE


async def vehicle_choice(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    query = update.callback_query
    vehicle_id = query.data
    context.user_data['vehicle_id'] = vehicle_id

    # Ask the user to enter the minimal count
    await query.edit_message_text(text=f'You selected: {vehicle_id}\n\nEnter the minimal count:')
    return MIN_COUNT_INPUT


async def min_count_input(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    min_count = update.message.text

    # Validate if the input is a number
    try:
        context.user_data['minimal_count'] = int(min_count)

    except ValueError:
        await update.message.reply_text('Please enter a valid number for the minimal count.')
        return MIN_COUNT_INPUT

    await STORAGE.save(
        storage.MonitoringData(vehicle_id=context.user_data['vehicle_id'], chat_id=update.message.chat_id,
                               minimal_count=context.user_data['minimal_count']))

    await update.message.reply_text(
        f'Success! Monitoring added!')
    return ConversationHandler.END


async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    await update.message.reply_text('Monitoring addition canceled.')
    return ConversationHandler.END


async def get_data_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('got getdata command. start processing')
    data = await POLLER.get_data()

    for entity in data:
        await context.bot.send_photo(update.message.chat_id, photo=entity.img, caption=str(entity))


async def get_data_short_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('got getshort command. start processing')
    data = await POLLER.get_data()

    await update.message.reply_text('\n\n'.join([str(v) for v in data]))


async def get_all_monitoring(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('got get_all_monitoring command. start processing')
    data = await STORAGE.get_all()
    message = '\n'.join([str(v) for v in data]) if data else 'No monitoring is set'
    await update.message.reply_text(message)


async def check_vehicles_count(context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('starting checking vehicles for min count')
    vehicles = await POLLER.get_data()
    for v in vehicles:
        users = [u for u in await STORAGE.get_by_vehicle_id(v.name) if v.current_count <= u.minimal_count]

        for u in users:
            await context.bot.send_photo(u.chat_id, v.img,
                                         caption=f'Attention! Only {v.current_count} {v.name}`s left. Current price is {v.price} gold')


def main() -> None:
    logging.basicConfig(
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.DEBUG
    )
    global POLLER, STORAGE
    POLLER = poller.Poller()
    STORAGE = storage.InMemoryStorage()

    app = ApplicationBuilder().token(BOT_TOKEN).build()

    vehicle_monitoring_handler = ConversationHandler(
        entry_points=[CommandHandler('add_monitoring', add_monitoring)],
        states={
            VEHICLE_CHOICE: [CallbackQueryHandler(vehicle_choice)],
            MIN_COUNT_INPUT: [MessageHandler(filters.TEXT & ~filters.COMMAND, min_count_input)],
        },
        fallbacks=[CommandHandler('cancel', cancel)],
        allow_reentry=True  # Allow users to restart the conversation by typing /add-monitoring again
    )

    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("data", get_data_command))
    app.add_handler(CommandHandler("data_short", get_data_short_command))
    app.add_handler(CommandHandler("monitoring", get_all_monitoring))
    app.add_handler(vehicle_monitoring_handler)

    def run_scheduler():
        while True:
            asyncio.run(check_vehicles_count(app))
            time.sleep(10)

    scheduler_thread = threading.Thread(target=run_scheduler)
    scheduler_thread.start()

    app.run_polling()


if __name__ == '__main__':
    main()
