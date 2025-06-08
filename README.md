# DeallsJobsTest API Documentation

This document provides comprehensive information about all the API endpoints available in the DeallsJobsTest application.

## Table of Contents

- [Authentication](#authentication)
- [Employee Endpoints](#employee-endpoints)
  - [Attendance](#attendance)
  - [Overtime](#overtime)
  - [Reimbursement](#reimbursement)
  - [Payslip](#payslip)
- [Admin Endpoints](#admin-endpoints)
  - [Period Management](#period-management)
  - [Payroll Processing](#payroll-processing)
  - [Payroll Summary](#payroll-summary)

## Authentication

### Login

Authenticates a user and returns a JWT token.

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Login successful",
      "token": "jwt_token_string"
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "message": "Invalid username or password",
      "error": "User not found"
    }
    ```

## Employee Endpoints

### Attendance

#### Submit Attendance

Submits attendance for the current user.

- **URL**: `/api/employee/attendance`
- **Method**: `POST`
- **Authentication**: Required (JWT Token)
- **Request Body**:
  ```json
  {
    "date": "YYYY-MM-DD HH:MM:SS"
  }
  ```
- **Description**: This is required to submit attendance the date should be the same date of today, if the date is not today it will return an error.
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Attendance submitted successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 409 Conflict | 403 Forbidden
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```

### Overtime

#### Submit Overtime

Submits overtime hours for the current user.

- **URL**: `/api/employee/overtime`
- **Method**: `POST`
- **Authentication**: Required (JWT Token)
- **Request Body**:
  ```json
  {
    "date": "YYYY-MM-DD HH:MM:SS",
    "hours": 1.5
  }
  ```
- **Success Response**:
  - **Code**: 201 Created
  - **Content**:
    ```json
    {
      "message": "Overtime submitted successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 403 Forbidden
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```

### Reimbursement

#### Submit Reimbursement

Submits a reimbursement request for the current user.

- **URL**: `/api/employee/reimbursement`
- **Method**: `POST`
- **Authentication**: Required (JWT Token)
- **Request Body**:
  ```json
  {
    "date": "YYYY-MM-DD",
    "amount": 100.50,
    "description": "Expense description"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Reimbursement submitted successfully",
      "ID": 123
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 500 Internal Server Error
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```

#### Upload Reimbursement Attachments

Uploads attachments for a reimbursement request.

- **URL**: `/api/employee/reimbursement/uploads/:id`
- **Method**: `POST`
- **Authentication**: Required (JWT Token)
- **URL Parameters**: `id` - ID of the reimbursement
- **Request Body**: Multipart form with `files` field containing the attachments
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Sucessfully uploaded files Reimbursement",
      "files": ["filename1.jpg", "filename2.pdf"]
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 403 Forbidden | 404 Not Found | 500 Internal Server Error
  - **Content**:
    ```json
    {
      "error": "Error message"
    }
    ```

### Payslip

#### Generate Payslip

Generates a payslip for the current user.

- **URL**: `/api/employee/payslip`
- **Method**: `GET`
- **Authentication**: Required (JWT Token)
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Payslip generated successfully",
      "data": {
        "employee_name": "John Doe",
        "period": "January 2023",
        "basic_salary": 5000.00,
        "overtime_pay": 150.00,
        "reimbursements": 200.00,
        "total_pay": 5350.00
      }
    }
    ```
- **Error Response**:
  - **Code**: 500 Internal Server Error
  - **Content**:
    ```json
    {
      "message": "Failed to generate payslip",
      "error": "Detailed error description"
    }
    ```

## Admin Endpoints

### Period Management

#### Create Period

Creates a new attendance period.

- **URL**: `/api/admin/period`
- **Method**: `POST`
- **Authentication**: Required (JWT Token with Admin privileges)
- **Request Body**:
  ```json
  {
    "start_date": "YYYY-MM-DD HH:MM:SS",
    "end_date": "YYYY-MM-DD HH:MM:SS",
    "is_locked": false
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Payroll retrieved successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 403 Forbidden | 500 Internal Server Error
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```

### Payroll Processing

#### Run Payroll

Initiates the payroll processing for a specific date.

- **URL**: `/api/admin/payroll`
- **Method**: `POST`
- **Authentication**: Required (JWT Token with Admin privileges)
- **Request Body**:
  ```json
  {
    "payroll_date": "YYYY-MM-DD HH:MM:SS"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Payroll running.."
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request | 403 Forbidden | 500 Internal Server Error
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```

### Payroll Summary

#### Get Payroll Summary

Retrieves the payroll summary data.

- **URL**: `/api/admin/summary/:period_id`
- **Method**: `GET`
- **Authentication**: Required (JWT Token with Admin privileges)
- **URL Parameters**: `period_id` - ID of the payroll period
- **Query Parameters**:
  - `page` (optional, default: 1): Page number for pagination
  - `limit` (optional, default: 10): Number of items per page
  - `sort` (optional): Sort order
  - `all` (optional, default: false): If true, returns all data without pagination
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Payroll retrieved successfully",
      "data": {
        "total_employees": 50,
        "total_payroll": 250000.00,
        "period": "January 2023",
        "payslips": [
          {
            "employee_id": 1,
            "employee_name": "John Doe",
            "total_pay": 5350.00
          },
          {
            "employee_id": 2,
            "employee_name": "Jane Smith",
            "total_pay": 4800.00
          }
        ]
      }
    }
    ```
- **Error Response**:
  - **Code**: 403 Forbidden | 404 Not Found | 500 Internal Server Error
  - **Content**:
    ```json
    {
      "message": "Error message",
      "error": "Detailed error description"
    }
    ```
