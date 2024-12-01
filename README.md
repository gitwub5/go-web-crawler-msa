# 📘 대학 공지사항 알림 시스템

이 프로젝트는 대학 공지사항을 크롤링하여 실시간으로 알람을 구독자에게 전달하거나, 웹 서버를 통해 사용자에게 제공하는 **마이크로서비스 아키텍처**로 설계되었습니다.

---

## **1. 아키텍처 개요**



---

## **2. 주요 컴포넌트**

### **1) 크롤링 서버 (Crawler Server)**

- **역할:**
    - 대학 웹사이트를 주기적으로 크롤링하여 새로운 공지사항 데이터를 감지.
    - 감지된 공지사항 데이터를 RabbitMQ에 메시지로 발행(Publisher).
- **연결:**
    - RabbitMQ 메시지 브로커에 메시지를 발행.

### **2) RabbitMQ**

- **역할:**
    - 크롤링 서버가 발행한 공지사항 메시지를 전달.
    - 웹 서버와 알람 서버가 RabbitMQ에 구독(Subscriber)하여 메시지를 수신.

### **3) 웹 서버 (Web Server)**

- **역할:**
    - RabbitMQ에서 공지사항 메시지를 구독.
    - MySQL 데이터베이스에 공지사항 데이터를 저장.
    - 사용자에게 공지사항 데이터를 제공(읽기 중심).

### **4) 알람 서버 (Alarm Server)**

- **역할:**
    - RabbitMQ에서 공지사항 메시지를 구독.
    - Redis를 이용해 실시간 알람 캐싱 및 전송.
    - MySQL을 이용해 구독자 데이터 관리 및 알림 전송 이력 저장.

### **5) MySQL**

- **역할:**
    - 공지 데이터와 구독자 데이터를 관리하기 위해 테이블 및 스키마 분리.
    - 크롤링 서버와 알람 서버에서 데이터를 저장 및 조회.

### **6) Redis**

- **역할:**
    - 알람 서버에서 실시간 알람 데이터를 캐싱.

---

## **3. 데이터 흐름**

### **1) 데이터 작성 흐름**

1. 크롤링 서버가 새로운 공지사항을 크롤링.
2. RabbitMQ에 공지사항 메시지를 발행.

### **2) 웹 서버 데이터 처리 흐름**

1. RabbitMQ에서 공지사항 메시지를 수신.
2. MySQL의 공지 테이블에 메시지 저장.
3. 사용자 요청 시 공지 데이터를 조회하여 반환.

### **3) 알람 서버 데이터 처리 흐름**

1. RabbitMQ에서 공지사항 메시지를 수신.
2. Redis에 공지 데이터를 캐싱.
3. 구독자 목록(MySQL)을 조회하여 푸시 알람 전송.
4. 알림 이력을 MySQL에 저장.

---

## **4. 데이터베이스 설계**

### **1) 공지 데이터 (웹 서버)**

- **Schema:** `crwl_db`

```sql
CREATE TABLE `cse` (
    `Number` VARCHAR(10) NOT NULL COMMENT '게시글 고유 번호',
    `Title` VARCHAR(200) NOT NULL COMMENT '글 제목',
    `Date` VARCHAR(20) NOT NULL COMMENT '등록 날짜',
    `Link` VARCHAR(200) NOT NULL COMMENT '링크',
    PRIMARY KEY (`Number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='컴퓨터공학과 공지사항';

CREATE TABLE `sw` (
    `Number` VARCHAR(10) NOT NULL COMMENT '게시글 고유 번호',
    `Title` VARCHAR(200) NOT NULL COMMENT '글 제목',
    `Date` VARCHAR(20) NOT NULL COMMENT '등록 날짜',
    `Link` VARCHAR(200) NOT NULL COMMENT '링크',
    PRIMARY KEY (`Number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='소프트웨어중심대학사업단 공지사항';
```

### **2) 구독자 데이터 및 알림 이력 (알람 서버)**

- **Schema:** `alarm_server_db`

#### **테이블 1:** `subscribers`

```sql
CREATE TABLE subscribers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    device_token VARCHAR(255) NOT NULL,
    subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### **테이블 2:** `notification_logs`

```sql
CREATE TABLE notification_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subscriber_id INT NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscriber_id) REFERENCES subscribers(id)
);
```

---

## **5. 확장성과 유지보수 전략**

### **RabbitMQ 사용 이유**

- 메시지 큐를 통해 웹 서버와 알람 서버를 분리.
- 대량의 메시지를 비동기적으로 전달하여 확장성 확보.

### **MySQL 분리 전략**

- 하나의 MySQL을 논리적으로 분리하여 사용.
- 트래픽 증가 시 물리적으로 데이터베이스를 분리 가능.

### **Redis 활용**

- 알람 서버에서 Redis를 사용해 실시간 알람 데이터 처리 성능 향상.