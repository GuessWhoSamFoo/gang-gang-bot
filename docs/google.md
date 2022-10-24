# Google Calendar API

Using the Calendar API requires a Google account with a project created.

There is [no cost](https://developers.google.com/calendar/api/guides/quota#pricing) for usage of the API.

## Setting up Google Calendar API

1. Go to [Google APIs & Services](https://console.cloud.google.com/) then go to `+ Enabled APIs & Services`
2. Enable Google Calendar API then go to credentials
3. Click `+ Create Credentials` and select `Service Account`
4. Under `Credential Type`, select Application data then choose `No` for using Google infrastructure.
5. Enter a service account name along with a description. Skip optional steps and click Done.
6. Click `Credentials`, then click the email of the newly created service account.
7. Under `Keys`, click `Add Key` and create a new JSON key.
8. Set the `GOOGLE_CREDENTIALS` environment variable to the contexts of the key json file.

## Initial Calendar Setup

1. Go to Google Calendar and click the three dots next to the calendar of interest
2. Under `Share with specific people`, add the service account email from the previous step with permission to make changes to events.
3. Under `Integrate calendar`, copy and paste the calendar ID to a config file or environment variable
