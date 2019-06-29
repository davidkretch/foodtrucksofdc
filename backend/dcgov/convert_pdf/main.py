import csv
import os

import camelot
from google.cloud import storage


def get_file(name, bucket, out):
    """get_file downloads a file from Google Cloud Storage.

    Args:
        name (str): The key of the file to download.
        bucket (str): The Google Cloud Storage bucket that stores the file.
        out (str): The folder to save the file to.
    
    Returns:
        str: The folder the file was saved to.
    """
    storage_client = storage.Client()
    b = storage_client.get_bucket(bucket)
    f = b.blob(name)
    path = os.path.join(out, os.path.basename(name))
    f.download_to_filename(path)
    return path


def convert_pdf_to_csv(name, out):
    """convert_pdf_to_csv converts a PDF containing a table to a CSV.

    Args:
        name (str): The path to the PDF file on disk.
        out (str): The folder to write the CSV to.
    
    Returns:
        str: The folder the file was saved to.
    """
    tables = camelot.read_pdf(name, pages='1-end', suppress_stdout=True)
    data = []
    for i, page in enumerate(tables):
        start = 1 if i == 0 else 2
        data.extend(page.data[start:])
    outname = os.path.splitext(os.path.basename(name))[0] + '.csv'
    path = os.path.join(out, outname)
    with open(path, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile, quoting=csv.QUOTE_NONNUMERIC)
        writer.writerows(data)
    return path


def save_file(name, bucket):
    """save_file saves a file to a GCS bucket.

    Args:
        name (str): The path to the file on disk.
        bucket (str): The name of the GCS bucket to write to.
    """
    storage_client = storage.Client()
    b = storage_client.get_bucket(bucket)
    f = b.blob(os.path.basename(name))
    f.upload_from_filename(name)


def convert_pdf(event, context):
    """convert_pdf converts a PDF to CSV and writes it back to Google Cloud Storage.
    
    Args:
        event (dict): Event payload.
        context (google.cloud.functions.Context): Metadata for the event.
    """
    bucket = event['bucket']
    name = event['name']
    print(f'Processing {name}')
    if os.path.splitext(name)[1] != '.pdf':
        return
    folder = '/tmp'
    pdf = get_file(name, bucket, folder)
    csv = convert_pdf_to_csv(pdf, folder)
    save_file(csv, bucket)
