CREATE TABLE [dbo].[rewards] (
    [reward_id] BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    [member_id] BIGINT NOT NULL,
    [point_type_id] INT NOT NULL,
    [points] BIGINT NOT NULL,
    [description] VARCHAR(500) NOT NULL,
    [event_date] DATETIME2(7) NOT NULL,

    CONSTRAINT FK_rewards_members
        FOREIGN KEY ([member_id])
        REFERENCES [dbo].[members]([member_id])
);