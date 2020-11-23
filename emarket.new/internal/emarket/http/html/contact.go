package html

const ContactTemplate = `
<h1 class="pageHeader text-center">{{.Title}}</h1>
<div class="contact-page">
    <h4 style="text-align: center;"><p>Связаться с нами можете любым удобным для Вас способом.</p><p>Мы рады Вам всегда! <i class="far fa-smile"></i></p></h4>
    <ul class="list-group contacts">
        <li class="list-group-item"><a class="nolink telegram" href="https://t.me/torg4u" target="_blank"><i class="fab fa-telegram-plane"></i> Telegram</a></li>
        <li class="list-group-item"><a class="nolink whatsapp" href="https://api.whatsapp.com/send?phone=79057740885" target="_blank">
            <i class="fab fa-whatsapp"></i> WhatsApp</a></li>
        <li class="list-group-item"><a class="nolink viber" href="viber://add?number=79057740885" target="_blank"><i class="fab fa-viber"></i> Viber</a></li>
        <li class="list-group-item"><a class="nolink email" href="mailto:support@torg4u.ru?subject=Вопрос" target="_blank"><i class="fas fa-at"></i> support@torg4u.ru</a></li>
        <li class="list-group-item"><a class="nolink phone" href="tel:+79057740885"><i class="fas fa-phone"></i> +7(905)774-08-85</a></li>
    </ul>
</div>
`
