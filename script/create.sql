-- POSTGRES SQL

--
-- value  
-- ├── ts           tag
-- ├── tag   ─────> ├── id                     sauk
-- ├── b            ├── sauk       ──┬───────> ├── id
-- ├── i            ├── sintep     ──│─┐       ├── name
-- └── r            ├── component  ──│─│─┐     └── comment
--                  ├── type         │ │ │      
--                  ├── signal       │ │ │     sintep
--                  └── name         │ ├─│───> ├── id
-- event                             │ │ │     └── name
-- ├── ts           msg              │ │ │
-- ├── msg   ─────> ├── id           │ │ │     component
-- ├── state ──┐    ├── sauk       ──┘ │ ├───> ├── id
-- └── text    │    ├── sintep     ────┘ │     ├── name
--             │    ├── component  ──────┘     └── descriptor
--             │    ├── severity
--             │    └── descriptor
--             │                               state
--             └─────────────────────────────> ├── id
--                                             ├── name
--                                             └── comment
--

---------------------------------------------------------------------------------------------------
-- sauk
-- В таблице хранится справочник САУК
--   id       цифровой уникальный не пустой идентификатор
--   name     уникальное не пустое короткое имя (без "САУК")
--   comment  полное расшифрованное наименование
CREATE TABLE "sauk"
(
  "id" bigint NOT NULL,
  "name" character varying(8) NOT NULL,
  "comment" character varying(80),
  CONSTRAINT "sauk_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "sauk_name_key" UNIQUE ("name"),
  CONSTRAINT "sauk_name_check" CHECK (COALESCE(TRIM(BOTH FROM "name"), ''::text) <> ''::text)
);
ALTER TABLE "sauk" OWNER to "user";
GRANT ALL ON TABLE "sauk" TO "user";

-- Заполнение данными таблицы sauk
-- Справочник должен быть одинаковым на всех объектах
INSERT INTO sauk("id","name","comment") VALUES
	(0,  '###',     'Нет данных'),
	(1,  'ВГП',     'Управление вентилятором главного проветривания'),
	(2,  'ВНУ',     'Управление воздухонагревательной установкой'),
	(3,  'ВУ',      'Управление водоотливной установкой'),
	(4,  'ГОУ',     'Управление газоотсасывающей установкой'),
	(5,  'ЗУ',      'Управление загрузочным устройством'),
	(6,  'К',       'Управление котельной'),
	(7,  'КТ',      'Управление конвейерным транспортом'),
	(8,  'МОНИТОР', 'Мониторинг персонала'),
	(9,  'ОС',      'Управление очистными сооружениями'),
	(10, 'ПМ',      'Управление подъемной машиной'),
	(11, 'ППНС',    'Управление противопожарной насосной станцией'),
	(12, 'ССС',     'Стволовая сигнализация и связь'),
	(13, 'УРП',     'Устройство регистрации параметров подъемной машины'),
	(14, 'Э',       'Управление электроснабжением');

---------------------------------------------------------------------------------------------------
-- sintep (системы)
--   id    цифровой уникальный не пустой идентификатор
--   name  уникальное не пустое имя
CREATE TABLE "sintep"
(
  "id" bigint NOT NULL,
  "name" character varying(80) NOT NULL,
  CONSTRAINT "sintep_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "sintep_name_key" UNIQUE ("name"),
  CONSTRAINT "sintep_name_check" CHECK (COALESCE(TRIM(BOTH FROM "name"), ''::text) <> ''::text)
);
ALTER TABLE "sintep" OWNER to "user";
GRANT ALL ON TABLE "sintep" TO "user";

-- Заполнение данными таблицы sintep
-- Справочник должен быть содеражать одинаковый код САУК на всех объектах
-- Справочник может содержать код единственной, собственной САУК 
-- Код должен формироваться по правилам формирования кода СИНТ (номер заказчика, номер системы ...)
INSERT INTO sintep("id","name") VALUES
	(0,'Нет такого объекта');

---------------------------------------------------------------------------------------------------
-- component (составная часть)
--   id       		цифровой уникальный не пустой идентификатор
--   name     		уникальное не пустое имя
--   descriptor		набор разнородных данных json для идентификации на основе описания
CREATE TABLE "component"
(
  "id" bigint NOT NULL,
  "name" character varying(25) NOT NULL,
  "descriptor" jsonb,
  CONSTRAINT "component_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "component_name_key" UNIQUE ("name"),
  CONSTRAINT "component_name_check" CHECK (COALESCE(TRIM(BOTH FROM "name"), ''::text) <> ''::text)
);
ALTER TABLE "component" OWNER to "user";
GRANT ALL ON TABLE "component" TO "user";
-- Индексы
CREATE INDEX "component_descriptor_index" ON "component" USING GIN ("descriptor");
ALTER INDEX "component_descriptor_index" OWNER to "user";

-- Заполнение данными таблицы component
-- Справочник должен быть одинаковым на всех объектах
INSERT INTO component("id","name") VALUES
	(0 , 'Нет такого компонента'),
	(1 , 'Общие сигналы'        ),
	(2 , 'ШУГ'                  ),
	(3 , 'ШУЛ-1'                ),
	(4 , 'ШУЛ-2'                ),
	(5 , 'ШУЛ-3'                ),
	(6 , 'ШУЛ-4'                ),
	(7 , 'ШУЛ-5'                ),
	(8 , 'ШУЛ-6'                ),
	(9 , 'ШУЛ-7'                ),
	(10, 'ШУЛ-8'                ),
	(11, 'ШУЛ-9'                ),
	(12, 'ШК'                   ),
	(13, 'ШИС'                  ),
	(14, 'ШПС'                  ),
	(15, 'Пульт оператора'      ),
	(16, 'Пульт диспетчера'     ),
	(17, 'Панель оператора'     ),
	(18, 'ПВШ'                  ),
	(19, 'ПВШ-1'                ),
	(20, 'ПВШ-2'                ),
	(21, 'ПВШ-3'                ),
	(22, 'ПВШ-4'                ),
	(23, 'ПВШ-5'                ),
	(24, 'ПВШ-6'                ),
	(25, 'ПВШ-7'                ),
	(26, 'ПВШ-8'                ),
	(27, 'ПВШ-9'                ),
	(28, 'МИБ'                  ),
	(29, 'МИБ-1'                ),
	(30, 'МИБ-2'                ),
	(31, 'МИБ-3'                ),
	(32, 'МИБ-4'                );

---------------------------------------------------------------------------------------------------
-- tag (теги scada системы)
-- Создание автонумератора для таблицы tag
CREATE SEQUENCE "tag_id_seq" CYCLE INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
ALTER SEQUENCE "tag_id_seq" OWNER TO "user";
GRANT ALL ON SEQUENCE "tag_id_seq" TO "user";
-- Создание таблицы tag 
--   id       	цифровой уникальный не пустой идентификатор
--   sauk	      идентификатор САУК (связан с таблицей sauk)
--   sintep     идентификатор объекта (связан с таблицей sintep)
--   component	идентификатор составной части (связан с таблицей component)
--   type       тип значений тега в символьном выражении "i" integer, "b" boolean, "r" real
--   signal     наименование обменного канала в scada системе
--   name       наименование тега в терминах пользователя
CREATE TABLE "tag"
(
  "id" bigint NOT NULL DEFAULT nextval('tag_id_seq'::regclass),
  "sauk" bigint,
  "sintep" bigint,
  "component" bigint,
  "type" char(1) NOT NULL DEFAULT 'i',
  "signal" character varying(80) NOT NULL,
  "name" character varying(80) NOT NULL,
  CONSTRAINT "tag_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "tag_type_check" CHECK (type = 'i' OR type = 'b' OR type = 'r'),
  CONSTRAINT "tag_sauk_id_ref" FOREIGN KEY ("sauk") REFERENCES "sauk" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "tag_sintep_id_ref" FOREIGN KEY ("sintep") REFERENCES "sintep" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "tag_component_id_ref" FOREIGN KEY ("component") REFERENCES "component" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "tag_signal_check" CHECK (COALESCE(TRIM(BOTH FROM "signal"), ''::text) <> ''::text)
);
ALTER TABLE "tag" OWNER to "user";
GRANT ALL ON TABLE "tag" TO "user";
-- Индексы
CREATE INDEX "tag_type_index" ON "tag" USING BTREE ("type");
ALTER INDEX "tag_type_index" OWNER to "user";
CREATE INDEX "tag_signal_index" ON "tag" USING BTREE ("signal");
ALTER INDEX "tag_signal_index" OWNER to "user";
CREATE INDEX "tag_name_index" ON "tag" USING BTREE ("name");
ALTER INDEX "tag_name_index" OWNER to "user";
-- Удалять автонумератор вместе с таблицей
ALTER SEQUENCE "tag_id_seq" OWNED BY "tag"."id";

---------------------------------------------------------------------------------------------------
-- value
-- В таблице хранятся значения тегов на момент времени
--   ts       	дата и время формирования значения
--   tag      	идентификатор тега (связан с таблицей tag)
--   b        	значение тип boolean
--   i        	значение тип integer
--   r        	значение тип real
CREATE TABLE "value"
(
  "ts" timestamp with time zone NOT NULL DEFAULT (now()),
  "tag" bigint NOT NULL,
  "b" boolean NOT NULL DEFAULT false,
  "i" integer NOT NULL DEFAULT 0,
  "r" real NOT NULL DEFAULT 0,
  CONSTRAINT "value_tag_id_ref" FOREIGN KEY ("tag") REFERENCES "tag" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID
);
ALTER TABLE "value" OWNER to "user";
GRANT ALL ON TABLE "value" TO "user";

-- ПОДСИСТЕМА АНАЛИЗА СОБЫТИЙ (см. ISA 18.2)
---------------------------------------------------------------------------------------------------
-- state 
-- В таблице хранится справочник состояний тревог
--   id       цифровой уникальный не пустой идентификатор
--   name     уникальное не пустое имя
--   comment  не обязательный комментарий
CREATE TABLE "state"
(
  "id" bigint NOT NULL,
  "name" character varying(5) NOT NULL,
  "comment" character varying(80),
  CONSTRAINT "state_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "state_name_key" UNIQUE ("name"),
  CONSTRAINT "state_name_check" CHECK (COALESCE(TRIM(BOTH FROM "name"), ''::text) <> ''::text)
);
ALTER TABLE "state" OWNER to "user";
GRANT ALL ON TABLE "state" TO "user";

-- Заполнение данными таблицы state
-- Справочник состояний тревог должен быть одинаковым на всех объектах
-- Список состояний в соответствии ISA 18.2
INSERT INTO state("id","name","comment") VALUES
	(0,'NUL','Нет данных'),
	(1,'NORM','Тревога не активна, подтверждена'),
	(2,'UNACK','Тревога активна, не подтверждена'),
	(3,'ACKED','Тревога активна, подтверждена'),
	(4,'RTNUN','Тревога не активна, не подтверждена'),
	(5,'SHLVD','Активация тревоги отложена оператором'),
	(6,'DSUPR','Активация тревоги запрещена программой'),
	(7,'OOSRV','Тревога выведена из эксплуатации'),
	(8,'ERR','Ошибка обработки');

---------------------------------------------------------------------------------------------------
-- msg
-- Создание автонумератора для таблицы msg
CREATE SEQUENCE "msg_id_seq" CYCLE INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
ALTER SEQUENCE "msg_id_seq" OWNER TO "user";
GRANT ALL ON SEQUENCE "msg_id_seq" TO "user";
-- Создание таблицы msg (сообщения scada системы)
--   id       	цифровой уникальный не пустой идентификатор
--   sauk		идентификатор САУК (связан с таблицей sauk)
--   sintep	идентификатор объекта (связан с таблицей sintep)
--   component	идентификатор составной части (связан с таблицей component)
--   severity	приоритет события
--   descriptor	набор данных json для идентификации на основе описания
--		рекомендуемые для реализации поля данных
--   	  dsupr		признак запрещено программой
--   	  oosrv		признак выведено из эксплуатации
--   	  defect   	признак причина брак оборудования
--   	  breakage	признак причина нарушение условий эксплуатации
--		  suspect		признак под подозрением
CREATE TABLE "msg"
(
  "id" bigint NOT NULL DEFAULT nextval('msg_id_seq'::regclass),
  "sauk" bigint,
  "sintep" bigint,
  "component" bigint,
  "severity" integer NOT NULL DEFAULT 100,
  "descriptor" jsonb,
  CONSTRAINT "msg_id_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "msg_sauk_id_ref" FOREIGN KEY ("sauk") REFERENCES "sauk" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "msg_sintep_id_ref" FOREIGN KEY ("sintep") REFERENCES "sintep" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "msg_component_id_ref" FOREIGN KEY ("component") REFERENCES "component" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID
);
ALTER TABLE "msg" OWNER to "user";
GRANT ALL ON TABLE "msg" TO "user";
-- Индексы
CREATE INDEX "msg_severity_index" ON "msg" USING BTREE ("severity");
ALTER INDEX "msg_severity_index" OWNER to "user";
CREATE INDEX "msg_descriptor_index" ON "msg" USING GIN ("descriptor");
ALTER INDEX "msg_descriptor_index" OWNER to "user";
-- Удалять автонумератор вместе с таблицей
ALTER SEQUENCE "msg_id_seq" OWNED BY "msg"."id";

---------------------------------------------------------------------------------------------------
-- event
-- Создание таблицы event
-- В таблице хранятся данные об событиях (событиях и тревогах)
--   ts       	дата и время события
--   msg      	идентификатор сообщения
--   state    	идентификатор состояния в которое перешла тревога (связан с таблицей state)
--   text       текст сообщения (на момент события)
-- События формируются алгоритмом, анализирующем значение как минимум одного тега
-- Текстовая строка сообщения может заполняться значениями через макрополя
CREATE TABLE "event"
(
  "ts" timestamp with time zone NOT NULL DEFAULT (now()),
  "msg" bigint NOT NULL,
  "state" bigint NOT NULL,
  "text" character varying(511) NOT NULL,
  CONSTRAINT "event_msg_id_ref" FOREIGN KEY ("msg") REFERENCES "msg" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "event_state_id_ref" FOREIGN KEY ("state") REFERENCES "state" ("id") MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID,
  CONSTRAINT "event_text_check" CHECK (COALESCE(TRIM(BOTH FROM "text"), ''::text) <> ''::text)
);
ALTER TABLE "event" OWNER to "user";
GRANT ALL ON TABLE "event" TO "user";
-- Индексы
CREATE INDEX "event_ts_index" ON "event" USING BTREE ("ts");
ALTER INDEX "event_ts_index" OWNER to "user";
