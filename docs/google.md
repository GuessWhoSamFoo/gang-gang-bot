## Setting up Google Calendar API

1. Go to [Google APIs & Services](https://console.cloud.google.com/) then go to `Enabled APIs & Services`
2. Enable Google Calendar API then go to credentials
3. Click `+ Create Credentials` and select `Service Account`
4. Fill in the form for the service account then click the email for the newly created account
5. Under `Keys`, `Add Key` then create a new key
6. Set the `GOOGLE_CREDENTIALS` environment variable to the contexts of the key json file.

## Initial Calendar Setup

1. Go to Google Calendar and click the three dots next to the calendar of interest
2. Under `Share with specific people`, add the service account email from the previous step with permission to make changes to events.
3. Under `Integrate calendar`, copy and paste the calendar ID to a config file or environment variable
