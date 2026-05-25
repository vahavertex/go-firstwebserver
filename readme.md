# Что конкретно изменено (Разбор ошибок и практик):
* Убран глобальный http.HandleFunc: Глобальный маршрутизатор (DefaultServeMux) опасен. Любой сторонний пакет, который вы импортируете, может тайно зарегистрировать в нем свой обработчик. Явный http.NewServeMux() решает эту проблему.
* Добавлены таймауты (ReadTimeout / WriteTimeout): Стандартный http.ListenAndServe(":8080", nil) не имеет ограничений по времени. Злоумышленник может открыть тысячи соединений и ничего в них не слать, что быстро перегрузит сервер.
* Замена fmt.Println на log.Println: Пакет log автоматически добавляет дату и время к каждой строчке, что критически важно для серверных логов.
* Замена panic(err) на log.Fatalf: Паника выводит весь стек вызовов (stack trace), что избыточно для ожидаемой ошибки запуска (например, если порт 8080 уже занят).
* Обновлен синтаксис путей (Go 1.22+): Строка "GET /{$}" гарантирует, что функция сработает только на главную страницу. Раньше обычный "/" перехватывал вообще любой запрос к серверу, из-за чего внутри функции приходилось писать проверку if r.URL.Path != "/" { http.NotFound(...) }.
Если вы хотите развить этот пример, мы можем добавить плавное завершение работы (Graceful Shutdown), чтобы сервер не обрывал активных клиентов при перезапуске, или настроить обработку JSON. Что из этого настроить?

[this is URL](https://markdown-editor.andona.click/#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)

≡ Краткое руководство Markdown

# Заголовок h1
## Заголовок h2
### Заголовок h3
#### Заголовок h4
##### Заголовок h5
###### Заголовок h6

Абзац Markdown. Пример:

Lorem ipsum dolor sit amet... Абзацы создаются при помощи пустой строки.

Для переноса строки делаем два пробела ` ` ` ` в конце (предыдущей) строки
Перенос строки

Получается? Отлично! :+1:

Текст с жирным начертанием (**strong**) и курсив (*italic*) в Markdown:

_1 символ_ `_` или `*` для наклонного текста
__2 символа__  `__` или `**` для жирного текста
***3 символа*** `___` или `***` для наклонного и жирного одновременно.

Перечеркнутый текст. 2 тильды `~` до и после текста - текст как перечеркнутый - ~~Зачеркнуто~~

Горизонтальная черта. `hr` - 3 звездочки или 3 дефиса

***

♦ Маркированный список. Для разметки неупорядоченных списков `*`, `-`, `+`:

* текст
* текст
* текст

Вложенные пункты. 4 пробела перед маркером:

* элемент маркированного списка
* элемент маркированного списка
    * вложенный текст
    * вложенный текст

Нумерованный список. Главное, чтобы перед элементом списка стояла цифра с точкой.

1. элемент нумерованного списка
2. элемент нумерованного списка
    1. вложенный
    2. вложенный

Можно сделать так:

0. текст
0. текст
0. текст

Список с абзацами:

* Текст
* Текст
* Текст

    Текст (4 пробела или `Tab`).

---

##### ♦ Ссылки Markdown

Здесь - [ссылка с title](https://example.com/ "Привет!").

Здесь - [ссылка без title](https://example.com/).

Ссылки с разметкой как у сносок.

Здесь - [ссылка][1] продолжение текста [ссылка][2] продолжение текста [ссылка][id]. [Просто ссылка][] без указания id.

[1]: https://example.com/ "Пример Title"
[2]: https://example.com/page
[id]: https://example.com/links (Пример Title)
[Просто ссылка]: https://example.com/short

Ссылки-сноски можно располагать в любом месте документа.

---

##### Цитаты в Markdown - cимвол `>`.

> Lorem ipsum dolor sit amet.
> Lorem ipsum dolor sit amet.
>
> Lorem ipsum dolor sit amet.

В цитаты можно помещать всё что угодно, в том числе вложенные цитаты:

> ### Заголовок.
>
> 1. список
> 2. список
>
> > Вложенная цитата.
>
> Исходный код (4 пробела в начале строки):
>
>     $source = file_get_contents('example.php');

##### Исходный код в Markdown

В GFM - поставить 3 апострофа (где `Ё`) до и после кода. Можно указать язык исходного кода.

```html
<div class="my-header">
    <h1>Матрёшка</h1>
    <p>Lorem ipsum dolor sit amet.</p>
</div>
```

```javascript
    $(function() { ... });
```

Для вставки кода внутри предложений - надо обрамить в апострофы (где `Ё`).

Пример: `<div class="my-markdown">`.

Если внутри кода есть апостроф, то код надо обрамить двойными апострофами: ``Бла-бла (`) тут.``

##### Картинки в Markdown

Картинка без alt текста

![](//placehold.co/200x100)

Картинка с alt и title:

![Alt text](//placehold.co/200x100 "Здесь title")

Картинка-ссылка:
Подсказка: синтаксис как у ссылок, только перед открывающей квадратной скобкой ставится восклицательный знак.

[![Alt text](//placehold.co/200x100)](https://example.com/)

Картинки-сноски:

![Картинка][image1]
![Картинка][image2]
![Картинка][image3]

[image1]: //placehold.co/200x100
[image2]: //placehold.co/150x100
[image3]: //placehold.co/100x100

---

##### Использование HTML внутри Markdown

Mожно смешивать Markdown и HTML. Если на какие-то элементы нужно поставить классы или атрибуты, используем HTML:

> Выделим слова без помощи * и _ . Например, это <em class="my-italic">курсив</em> и это тоже <i>курсив</i>. А вот так уже <b>strong</b>, и так тоже <strong>strong</strong>.

Можно и наоборот, внутри HTML-тегов использовать Markdown.

<div class="my-markdown">

###### Markdown внутри HTML. Пример:

Выделять слова можно при помощи `*` и `_` . Например, это _курсив_ и это тоже *italic*. А вот так уже __strong__, и так тоже **strong**.

</div>

---

##### Таблицы

В чистом Маркдауне нет синтаксиса для таблиц, а в GFM есть. Рисуем:

First Header  | Second Header
------------- | -------------
Content Cell  | Content Cell
Content Cell  | Content Cell

Можно по бокам линии нарисовать:

| First Header  | Second Header |
| ------------- | ------------- |
| Content Cell  | Content Cell  |
| Content Cell  | Content Cell  |

Можно управлять выравниванием столбцов при помощи двоеточия:

| Left-Aligned  | Center Aligned  | Right Aligned |
|:------------- |:---------------:| -------------:|
| col 3 is      | какой-то текст  |   **my text** |
| col 2 is      | центр           |           $123|
| Content Cell  | бука            |         ~~$7~~|

Внутри таблиц можно использовать ссылки, наклонный, жирный или зачеркнутый текст.

---

♦ Для всего остального есть обычный HTML.

---

###### Links:

 * <small>[markdown-it](https://github.com/markdown-it/markdown-it) for Markdown parsing</small>
 * <small>[CodeMirror](https://codemirror.net/) for the awesome syntax-highlighted editor</small>
 * <small>[Live (Github-flavored)](https://github.com/jbt/markdown-editor) Markdown Editor</small>
 * <small>[highlight.js](https://softwaremaniacs.org/soft/highlight/en/) for syntax highlighting in output code blocks</small>
 * <small>[js-deflate](https://github.com/dankogai/js-deflate) for gzipping of data to make it fit in URLs</small>
