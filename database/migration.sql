DROP database IF EXISTS exercise;
Go

Create database exercise;
Go

USE exercise;
Go

Create table company (
    uuid uniqueidentifier not null CONSTRAINT DF_company_uuid DEFAULT ((newid())),
    company_name nvarchar(255),
    code nvarchar(255),
    country nvarchar(128),
    website nvarchar(255),
    phone nvarchar(32),
    archive bit not null CONSTRAINT DF_company_archive DEFAULT ((0)),
    dt_created datetime NOT NULL CONSTRAINT DF_company_dt_created DEFAULT (GETUTCDATE()),
    dt_updated datetime NOT NULL CONSTRAINT DF_company_dt_updated DEFAULT (GETUTCDATE()),
    dt_archived datetime NULL
);

Insert into Company (company_name, code, country, website, phone) values
('FirstTestCompanyName','FirstTestCompanyCode','FirstTestCompanyCountry','FirstTestCompanyWebsite','FirstTestCompanyPhone'),
('SecondTestCompanyName','SecondTestCompanyCode','SecondTestCompanyCountry','SecondTestCompanyWebsite','SecondTestCompanyPhone');
