#!/usr/bin/env bash

sdp="eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz1tb3ppbGxhLi4uVEhJU19JU19TRFBBUlRBLTg2LjAuMSA3MTA4ODQ3Njc1OTk2MDYwODAxIDAgSU4gSVA0IDAuMC4wLjBcclxucz0tXHJcbnQ9MCAwXHJcbmE9c2VuZHJlY3ZcclxuYT1maW5nZXJwcmludDpzaGEtMjU2IDU1OkRBOjk3OkJBOjRBOkNFOkFFOjY2OjFFOjU1OkRGOkM1OjNFOjMwOjlCOkU3OjIxOjkxOjhEOjY1OjA1Ojg1OjlGOkYyOjZFOjRGOkUxOjBDOkVGOjI2OkE5OjUyXHJcbmE9Z3JvdXA6QlVORExFIDBcclxuYT1pY2Utb3B0aW9uczp0cmlja2xlXHJcbmE9bXNpZC1zZW1hbnRpYzpXTVMgKlxyXG5tPXZpZGVvIDkgVURQL1RMUy9SVFAvU0FWUEYgMTIwIDEyNCAxMjEgMTI1IDEyNiAxMjcgOTcgOThcclxuYz1JTiBJUDQgMC4wLjAuMFxyXG5hPXNlbmRyZWN2XHJcbmE9ZW5kLW9mLWNhbmRpZGF0ZXNcclxuYT1leHRtYXA6MyB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOm1pZFxyXG5hPWV4dG1hcDo0IGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L2Ficy1zZW5kLXRpbWVcclxuYT1leHRtYXA6NSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDp0b2Zmc2V0XHJcbmE9ZXh0bWFwOjYvcmVjdm9ubHkgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvcGxheW91dC1kZWxheVxyXG5hPWV4dG1hcDo3IGh0dHA6Ly93d3cuaWV0Zi5vcmcvaWQvZHJhZnQtaG9sbWVyLXJtY2F0LXRyYW5zcG9ydC13aWRlLWNjLWV4dGVuc2lvbnMtMDFcclxuYT1mbXRwOjEyNiBwcm9maWxlLWxldmVsLWlkPTQyZTAxZjtsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0xXHJcbmE9Zm10cDo5NyBwcm9maWxlLWxldmVsLWlkPTQyZTAxZjtsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xXHJcbmE9Zm10cDoxMjAgbWF4LWZzPTEyMjg4O21heC1mcj02MFxyXG5hPWZtdHA6MTI0IGFwdD0xMjBcclxuYT1mbXRwOjEyMSBtYXgtZnM9MTIyODg7bWF4LWZyPTYwXHJcbmE9Zm10cDoxMjUgYXB0PTEyMVxyXG5hPWZtdHA6MTI3IGFwdD0xMjZcclxuYT1mbXRwOjk4IGFwdD05N1xyXG5hPWljZS1wd2Q6YmIwYzNkNTIxZmZkODRkNTQ0ZDk2MTM5ODQ2OGNmNGJcclxuYT1pY2UtdWZyYWc6YTc4MzdjZmRcclxuYT1taWQ6MFxyXG5hPW1zaWQ6LSB7YTVhMWJmZmItYzM1Ni0xMjRjLTkxYzYtNjY5ZGRmMDE2ZGE1fVxyXG5hPXJ0Y3AtZmI6MTIwIG5hY2tcclxuYT1ydGNwLWZiOjEyMCBuYWNrIHBsaVxyXG5hPXJ0Y3AtZmI6MTIwIGNjbSBmaXJcclxuYT1ydGNwLWZiOjEyMCBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEyMCB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEyMSBuYWNrXHJcbmE9cnRjcC1mYjoxMjEgbmFjayBwbGlcclxuYT1ydGNwLWZiOjEyMSBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMjEgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMjEgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMjYgbmFja1xyXG5hPXJ0Y3AtZmI6MTI2IG5hY2sgcGxpXHJcbmE9cnRjcC1mYjoxMjYgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTI2IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTI2IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6OTcgbmFja1xyXG5hPXJ0Y3AtZmI6OTcgbmFjayBwbGlcclxuYT1ydGNwLWZiOjk3IGNjbSBmaXJcclxuYT1ydGNwLWZiOjk3IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6OTcgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1tdXhcclxuYT1ydGNwLXJzaXplXHJcbmE9cnRwbWFwOjEyMCBWUDgvOTAwMDBcclxuYT1ydHBtYXA6MTI0IHJ0eC85MDAwMFxyXG5hPXJ0cG1hcDoxMjEgVlA5LzkwMDAwXHJcbmE9cnRwbWFwOjEyNSBydHgvOTAwMDBcclxuYT1ydHBtYXA6MTI2IEgyNjQvOTAwMDBcclxuYT1ydHBtYXA6MTI3IHJ0eC85MDAwMFxyXG5hPXJ0cG1hcDo5NyBIMjY0LzkwMDAwXHJcbmE9cnRwbWFwOjk4IHJ0eC85MDAwMFxyXG5hPXNldHVwOmFjdHBhc3NcclxuYT1zc3JjOjE2NDg4ODA1NDEgY25hbWU6e2E2ZmRlNGExLTJkYWUtOTg0Ny04ZDE3LWYzM2QyNDllZGU0MH1cclxuIn0="

curl localhost:8080/sdp -d "${sdp}"
