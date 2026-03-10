
# DragonKi – Use Case
## Automatic Daily Backup with External Storage

---

# 1. Overview

To ensure the safety and integrity of player progress and the game economy, DragonKi implements an **automatic daily backup system** of the production database.

Backups are exported from the Railway database and securely stored in an **external storage provider** (e.g., Cloudflare R2, AWS S3, or compatible object storage).

This mechanism guarantees that game data can be restored in case of system failure, database corruption, or operational mistakes.

---

# 2. Actors

| Actor | Description |
|------|-------------|
| Backup Service | Automated worker responsible for generating and uploading backups |
| Railway Database | Primary PostgreSQL database storing game data |
| External Storage | Remote storage used to persist backup files |
| Game Administrator | Responsible for restoring backups if necessary |

---

# 3. Preconditions

- The DragonKi database is running on Railway.
- Environment variable `DATABASE_URL` is available.
- External storage credentials are configured.
- Backup service has permission to access the database.

---

# 4. Backup Frequency

| Type | Frequency |
|-----|-----------|
| Full Backup | Daily |
| Retention | 30 days |
| Compression | Enabled |

Backups are scheduled using a **cron job or worker service**.

Example schedule:

Every day at 03:00 UTC

---

# 5. Backup Process Flow

1. Backup service is triggered by scheduled job.
2. The service connects to the Railway PostgreSQL database.
3. A full database dump is generated using `pg_dump`.
4. The dump file is compressed.
5. The compressed file is uploaded to external storage.
6. Backup retention policy removes old backups.

---

# 6. Technical Flow

Railway PostgreSQL  
↓  
Backup Worker (cron job)  
↓  
pg_dump execution  
↓  
Backup compression (.gz)  
↓  
Upload to External Storage  

---

# 7. Backup File Naming

Backup files follow a standardized naming convention:

dragonki-backup-YYYY-MM-DD.sql.gz

Example:

dragonki-backup-2026-03-09.sql.gz

---

# 8. Storage Structure

Example structure in external storage:

dragonki-backups/
    2026/
        03/
            dragonki-backup-2026-03-09.sql.gz
            dragonki-backup-2026-03-08.sql.gz
            dragonki-backup-2026-03-07.sql.gz

---

# 9. Critical Data Protected

The following tables are considered critical and must always be included in backups:

- players
- inventory
- equipment
- battles
- transactions
- economy
- big_boss_damage
- progression

Loss of these tables may result in permanent loss of player progress.

---

# 10. Restore Procedure

In case of failure, the administrator may restore a backup using:

gunzip dragonki-backup.sql.gz  
psql DATABASE_URL < dragonki-backup.sql

Steps:

1. Download backup from external storage.
2. Decompress the file.
3. Restore database using `psql`.

---

# 11. Security Considerations

- Backup files must be stored in **private storage buckets**.
- Access keys must be stored as environment variables.
- Backup service must use encrypted connections.

---

# 12. Benefits

This backup architecture ensures:

- Protection against database corruption
- Recovery from operational mistakes
- Long-term preservation of player progress
- Independent storage outside Railway infrastructure

---

# 13. Future Improvements

Potential improvements include:

- Incremental backups
- Real-time replication
- Automated restore testing
- Multi-region storage redundancy

---

**DragonKi – Data Protection Architecture**
