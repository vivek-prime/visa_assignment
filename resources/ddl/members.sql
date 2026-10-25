CREATE TABLE [dbo].[members] (
    [member_id] BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    [name] VARCHAR(100) NOT NULL,
    [email] VARCHAR(255) NOT NULL,
    [created_at] DATETIME2(7) NOT NULL
);