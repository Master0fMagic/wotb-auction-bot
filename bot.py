import logging

from telegram.ext import ApplicationBuilder, CommandHandler, ContextTypes
from telegram import Update
import poller

BOT_TOKEN = ''
POLLER: poller.Poller


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text('Hello! I am your bot. Use /getdata to fetch the latest data.')


async def get_data_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('got getdata command. start processing')
    data = await POLLER.get_data()

    for entity in data:
        await context.bot.send_photo(update.message.chat_id, photo=entity.img, caption=str(entity))


async def get_data_short_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    logging.debug('got getshort command. start processing')
    data = await POLLER.get_data()

    await update.message.reply_text('\n\n'.join([str(v) for v in data]))


def main() -> None:
    logging.basicConfig(
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.DEBUG
    )
    global POLLER
    POLLER = poller.Poller()

    app = ApplicationBuilder().token(BOT_TOKEN).build()

    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("data", get_data_command))
    app.add_handler(CommandHandler("datashort", get_data_short_command))

    app.run_polling()


if __name__ == '__main__':
    main()
