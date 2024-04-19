DELETE FROM "sintep";

INSERT INTO sintep("id","name") VALUES
	(18014,'ш. Кирова, Водоотлив центральный');

INSERT INTO tag("sauk","sintep","component","type","signal","name") VALUES
	(3,18014,1,'b','com.bit_00_00','Задвижка М1.1 Аварийный стоп'),
	(3,18014,1,'b','com.bit_00_01','Задвижка М1.1 Авария'),
	(3,18014,1,'b','com.bit_00_02','Задвижка М1.1 Момент при закрытии больше аварийного'),
	(3,18014,1,'b','com.bit_00_03','Задвижка М1.1 Момент при открытии больше аварийного'),
	(3,18014,1,'b','com.bit_00_04','Задвижка М1.1 Не закрылась'),
	(3,18014,1,'b','com.bit_00_05','Задвижка М1.1 Не открылась'),
	(3,18014,1,'b','com.bit_00_06','Задвижка М1.1 Нет готовности'),
	(3,18014,1,'b','com.bit_00_07','Задвижка М1.1 Отказ блока управления'),
	(3,18014,1,'b','com.bit_00_08','Насос 1. Аварийный стоп'),
	(3,18014,1,'b','com.bit_00_09','Насос 1. Авария'),
	(3,18014,1,'b','com.bit_00_10','Насос 1. Время запуска более 10 минут'),
	(3,18014,1,'b','com.bit_00_11','Насос 1. Давление выше аварийного'),
	(3,18014,1,'b','com.bit_00_12','Насос 1. Давление на всасе выше аварийного'),
	(3,18014,1,'b','com.bit_00_13','Насос 1. Давление на всасе ниже аварийного'),
	(3,18014,1,'b','com.bit_00_14','Насос 1. Давление на напоре не достигло требуемого'),
	(3,18014,1,'b','com.bit_00_15','Насос 1. Давление на напоре снизилось до минимума'),
	(3,18014,1,'b','com.bit_01_00','Насос 1. Давление ниже аварийного'),
	(3,18014,1,'b','com.bit_01_01','Насос 1. Задвижка на всасе самопроизвольно закрылась'),
	(3,18014,1,'b','com.bit_01_02','Насос 1. Задвижка на напоре самопроизвольно закрылась'),
	(3,18014,1,'b','com.bit_01_03','Насос 1. КМШ. Авария по току'),
	(3,18014,1,'b','com.bit_01_04','Насос 1. КМШ. Нет связи с ячейкой'),
	(3,18014,1,'b','com.bit_01_05','Насос 1. Не включился'),
	(3,18014,1,'b','com.bit_01_06','Насос 1. Не отключился'),
	(3,18014,1,'b','com.bit_01_07','Насос 1. Подшипник 1. Виброскорость выше аварийной'),
	(3,18014,1,'b','com.bit_01_08','Насос 1. Подшипник 1. Температура выше аварийной'),
	(3,18014,1,'b','com.bit_01_09','Насос 1. Подшипник 2. Виброскорость выше аварийной'),
	(3,18014,1,'b','com.bit_01_10','Насос 1. Подшипник 2. Температура выше аварийной'),
	(3,18014,1,'b','com.bit_01_11','Насос 1. Работа на закрытую задвижку'),
	(3,18014,1,'b','com.bit_01_12','Насос 1. Расход выше аварийного'),
	(3,18014,1,'b','com.bit_01_13','Насос 1. Расход ниже аварийного'),
	(3,18014,1,'b','com.bit_01_14','Насос 1. Самопроизвольное отключение'),
	(3,18014,1,'b','com.bit_01_15','Насос 1. Самопроизвольный запуск'),
	(3,18014,1,'b','com.bit_02_00','Нет работающих насосов'),
	(3,18014,1,'b','com.bit_02_01','Общий трубопровод 1. Давление выше аварийного'),
	(3,18014,1,'b','com.bit_02_02','Общий трубопровод 1. Давление ниже аварийного'),
	(3,18014,1,'b','com.bit_02_03','Общий трубопровод 2. Давление выше аварийного'),
	(3,18014,1,'b','com.bit_02_04','Общий трубопровод 2. Давление ниже аварийного'),
	(3,18014,1,'b','com.bit_02_05','Общий трубопровод 3. Давление выше аварийного'),
	(3,18014,1,'b','com.bit_02_06','Общий трубопровод 3. Давление ниже аварийного'),
	(3,18014,1,'b','com.real_00','Двигатель 1. Подшипник 1. Температура'),
	(3,18014,1,'b','com.real_01','Двигатель 1. Подшипник 2. Температура'),
	(3,18014,1,'b','com.real_02','Насос 1. Виброскорость'),
	(3,18014,1,'b','com.real_03','Насос 1. Время заливки'),
	(3,18014,1,'b','com.real_04','Насос 1. Давление'),
	(3,18014,1,'b','com.real_05','Насос 1. Подшипник 1. Виброскорость'),
	(3,18014,1,'b','com.real_06','Насос 1. Подшипник 1. Температура'),
	(3,18014,1,'b','com.real_07','Насос 1. Подшипник 2. Виброскорость'),
	(3,18014,1,'b','com.real_08','Насос 1. Подшипник 2. Температура');

INSERT INTO msg("sauk","sintep","component","severity","descriptor") VALUES
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Аварийный стоп"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Авария"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Момент при закрытии больше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Момент при открытии больше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Не закрылась"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Не открылась"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Нет готовности"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Задвижка М1.1 Отказ блока управления"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Аварийный стоп"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Авария"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Время запуска более 10 минут"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление на всасе выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление на всасе ниже аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление на напоре не достигло требуемого"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление на напоре снизилось до минимума"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Давление ниже аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Задвижка на всасе самопроизвольно закрылась"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Задвижка на напоре самопроизвольно закрылась"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. КМШ. Авария по току"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. КМШ. Нет связи с ячейкой"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Не включился"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Не отключился"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Подшипник 1. Виброскорость выше аварийной"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Подшипник 1. Температура выше аварийной"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Подшипник 2. Виброскорость выше аварийной"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Подшипник 2. Температура выше аварийной"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Работа на закрытую задвижку"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Расход выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Расход ниже аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Самопроизвольное отключение"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Насос 1. Самопроизвольный запуск"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Нет работающих насосов"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 1. Давление выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 1. Давление ниже аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 2. Давление выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 2. Давление ниже аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 3. Давление выше аварийного"}'),
	(3,18014,1,600,'{"dsupr":"false","oosrv":"false","text":"Общий трубопровод 3. Давление ниже аварийного"}');
