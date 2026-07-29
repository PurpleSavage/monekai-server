# API Endpoints Documentation

Base URL: `http://localhost:8080`

---

## Authentication (`/auth`)

### `POST /auth/login`
Inicia sesión con Google OAuth.

- **Auth:** Pública
- **Body (JSON):**
  ```json
  {
    "token": "string (requerido) — token de Google OAuth"
  }
  ```
- **Response (200):** `ResponseSessionDto`
  ```json
  {
    "userData": {
      "id": "uuid",
      "email": "string",
      "photoUrl": "string | null",
      "createdAt": "string (ISO 8601)",
      "credits": "int"
    },
    "accessToken": "string (JWT, expira en 1 minuto)"
  }
  ```
- **Cookie:** Setea `session_token` (refresh token JWT, expira en 72h)

---

### `GET /auth/refresh-token`
Refresca el access token usando el refresh token.

- **Auth:** Cookie `session_token` (RefreshToken)
- **Response (200):** `RefreshTokenResponseDto`
  ```json
  {
    "accessToken": "string (nuevo JWT)"
  }
  ```

---

### `GET /auth/profile`
Obtiene el perfil del usuario autenticado y renueva la sesión.

- **Auth:** Cookie `session_token` (RefreshToken)
- **Response (200):** `RenovateSessionResponseDto`
  ```json
  {
    "userData": {
      "id": "uuid",
      "email": "string",
      "photoUrl": "string | null",
      "createdAt": "string (ISO 8601)",
      "credits": "int"
    },
    "accessToken": "string (nuevo JWT)"
  }
  ```

---

## Sampler / Audio (`/audio`)

### `POST /audio/create`
Crea una solicitud de generación de audio vía Replicate.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Middleware extra:** `CheckCredits` (verifica que el usuario tenga créditos suficientes)
- **Body (JSON):**
  ```json
  {
    "prompt": "string (requerido, min 5, max 500 caracteres)",
    "modelVersion": "string (requerido) — stereo-melody-large | stereo-large | melody-large | large",
    "duration": "int (requerido, 1-30 segundos)",
    "outputFormat": "string (requerido) — mp3 | wav",
    "sampleName": "string (requerido, min 1, max 500)",
    "email": "string (requerido, email válido)"
  }
  ```
- **Response (200):** Objeto de la transacción de Replicate (respuesta inmediata, el procesamiento es asíncrono)

---

### `GET /audio/samples`
Lista los samples del usuario autenticado (con paginación).

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Query Params:**
  - `page` (int, requerido)
  - `limit` (int, requerido)
- **Response (200):** `PaginatedResponse<SampleResponseDTO>`
  ```json
  {
    "total": "int",
    "limit": "int",
    "page": "int",
    "data": [
      {
        "id": "uuid",
        "sampleName": "string",
        "prompt": "string",
        "audioUrl": "string | null",
        "duration": "int",
        "outputFormat": "mp3 | wav",
        "modelVersion": "string",
        "status": "starting | processing | succeeded | failed | canceled",
        "createdAt": "string"
      }
    ]
  }
  ```

---

### `POST /audio/share-sample`
Comparte un sample en la comunidad.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Body (JSON):**
  ```json
  {
    "sampleId": "uuid (requerido)"
  }
  ```
- **Response (200):** `ShareSampleResponseDTO`
  ```json
  {
    "id": "uuid",
    "sampleId": "uuid | null",
    "sampleVersionId": "uuid | null",
    "userId": "string",
    "likes": 0,
    "downloads": 0,
    "createdAt": "string"
  }
  ```

---

### `POST /audio/sample-version`
(Pendiente de implementación)

---

### `POST /audio/webhook/songs`
Webhook de Replicate. Recibe la respuesta de la generación de audio.

- **Auth:** Middleware `ReplicateMiddlewareWebhook` (verifica firma del webhook)
- **Body:** Payload de Replicate (JSON)
- **Response (200):** Vacío (procesa internamente)
- **Efectos secundarios:**
  - Actualiza la URL del sample generado
  - Dispara evento `sample_event` vía Observer → notificación SSE al usuario

---

## Community (`/community`)

### `GET /community/samples`
Lista samples compartidos (samples originales).

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Query Params:**
  - `page` (int, requerido)
  - `limit` (int, requerido)
- **Response (200):** Array de `SharedSampleItemDTO`
  ```json
  [
    {
      "id": "uuid",
      "likes": "int",
      "downloads": "int",
      "createdAt": "string",
      "sample": {
        "id": "uuid",
        "sampleName": "string",
        "initialAudioUrl": "string",
        "prompt": "string",
        "duration": "int"
      },
      "sharedBy": {
        "userId": "uuid",
        "name": "string",
        "email": "string"
      }
    }
  ]
  ```

---

### `GET /community/edit-samples`
Lista samples compartidos con versiones editadas.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Query Params:**
  - `page` (int, requerido)
  - `limit` (int, requerido)
- **Response (200):** Array de `SharedSampleVersionItemDTO`
  ```json
  [
    {
      "id": "uuid",
      "likes": "int",
      "downloads": "int",
      "createdAt": "string",
      "sampleVersion": {
        "id": "uuid",
        "effects": { ... },
        "finalAudioUrl": "string",
        "sampleName": "string",
        "prompt": "string"
      },
      "sharedBy": {
        "userId": "uuid",
        "name": "string",
        "email": "string"
      }
    }
  ]
  ```

---

### `PATCH /community/like/{sampleID}`
Da like a un sample compartido.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Path Param:**
  - `sampleID` (uuid, requerido)
- **Response (200):** `LikeSharedSampleResponseDTO`
  ```json
  {
    "message": "string",
    "sampleIDModify": "uuid"
  }
  ```

---

### `GET /community/download/{sampleID}`
Descarga un sample compartido (incrementa contador de descargas).

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Path Param:**
  - `sampleID` (uuid, requerido)
- **Response (200):** `DownloadSharedSampleVO`
  ```json
  {
    "sampleID": "uuid",
    "downloads": "int"
  }
  ```

---

## Payments (`/payments`)

### `POST /payments/create`
Crea una transacción de pago en Paddle.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Query Params:**
  - `packageId` (uuid, requerido — ID del credit package)
- **Response (200):** `CreateTransactionResponseDTO`
  ```json
  {
    "TransactionID": "string (txn_xxxxxx)",
    "Status": "draft | ready | billed | paid | completed | canceled | past_due"
  }
  ```

---

### `POST /payments/webhook`
Webhook de Paddle. Recibe eventos de transacciones.

- **Auth:** Middleware `PaymentWebhookVerifier` (verifica firma HMAC de Paddle)
- **Body:** Payload de Paddle (JSON, verificado automáticamente)
- **Response (200):** Vacío
- **Efectos secundarios:**
  - Procesa eventos `transaction.billed` y `transaction.completed`
  - Guarda el pago en BD (con idempotencia por `provider_transaction_id`)
  - Suma créditos al usuario (`UPDATE users SET credits = credits + N`)
  - Dispara evento `payment_event` vía Observer → notificación SSE al usuario

---

## Notifications (`/notifications`)

### `GET /notifications/all`
Lista las notificaciones del usuario autenticado (con paginación).

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Query Params:**
  - `page` (int, requerido, min 1)
  - `limit` (int, requerido, min 1, max 100)
- **Response (200):** `PaginatedResponse<ItemNotificationDTO>`
  ```json
  {
    "total": "int",
    "limit": "int",
    "page": "int",
    "data": [
      {
        "ID": "uuid",
        "Type": "replicate_error | replicate_success | payment | info",
        "Title": "string",
        "Message": "string",
        "Status": "unread | read",
        "ReferenceID": "uuid",
        "CreatedAt": "string",
        "UserID": "uuid"
      }
    ]
  }
  ```

---

### `PATCH /notifications/read-all`
Marca múltiples notificaciones como leídas.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Body (JSON):**
  ```json
  {
    "notificationIds": ["uuid", "uuid", ...]
  }
  ```
- **Response (200):** `NotificationMarkResponseDTO`
  ```json
  {
    "notificationId": "uuid",
    "statusNotification": "read"
  }
  ```

---

### `PATCH /notifications/{id}/read`
Marca una notificación como leída.

- **Auth:** Header `Authorization: Bearer <accessToken>` (AccessToken)
- **Path Param:**
  - `id` (uuid, requerido — ID de la notificación)
- **Body (JSON):**
  ```json
  {
    "notificationId": "uuid"
  }
  ```
- **Response (200):** `NotificationMarkResponseDTO`
  ```json
  {
    "notificationId": "uuid",
    "statusNotification": "read"
  }
  ```

---

## Server-Sent Events (`/sse`)

### `GET /sse/stream`
Conexión SSE para recibir eventos en tiempo real.

- **Auth:** Cookie `session_token` (RefreshToken)
- **Headers de respuesta:**
  ```
  Content-Type: text/event-stream
  Cache-Control: no-cache
  Connection: keep-alive
  ```
- **Formato de eventos:** (SSE estándar)
  ```
  event: sample_ready
  data: {"id":"...","type":"replicate_success","title":"Sample ready","message":"...","status":"unread","reference_id":"...","created_at":"...","user_id":"...","email":"...","data":{...}}
  
  event: sample_error
  data: {"id":"...","type":"replicate_error","title":"Sample failed","message":"...","status":"unread","reference_id":"...","created_at":"...","user_id":"...","email":"...","data":null}
  
  event: payment_success
  data: {"id":"...","type":"payment","title":"Payment successful","message":"Credits have been added to your account","status":"unread","reference_id":"...","created_at":"...","user_id":"...","email":"...","data":{"payment_id":"...","credits":100,"amount":2500,"currency":"USD","status":"completed"}}
  
  event: payment_failed
  data: {"id":"...","type":"payment","title":"Payment failed","message":"The payment could not be processed","status":"unread","reference_id":"...","created_at":"...","user_id":"...","email":"...","data":{"payment_id":"...","credits":0,"amount":0,"currency":"","status":""}}
  ```
- **Keep-alive:** Cada 15 segundos se envía un comentario vacío (`: keep-alive`)

---

## Middleware Reference

| Middleware | Tipo | Módulo | Endpoints |
|-----------|------|--------|-----------|
| `AccessToken` | Header `Authorization: Bearer <JWT>` | Auth (shared) | `/audio/*`, `/community/*`, `/payments/create`, `/notifications/*` |
| `RefreshToken` | Cookie `session_token` | Auth (shared) | `/auth/refresh-token`, `/auth/profile`, `/sse/stream` |
| `CheckCredits` | Header `Authorization` + verifica créditos | Shared | `/audio/create` |
| `ReplicateMiddlewareWebhook` | Verifica firma HMAC de Replicate | Sampler | `/audio/webhook/songs` |
| `PaymentWebhookVerifier` | Verifica firma HMAC de Paddle | Payments | `/payments/webhook` |
